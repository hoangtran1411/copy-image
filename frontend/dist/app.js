/**
 * Copy Image Tool - Frontend Application
 * 
 * This file handles the UI interactions and communicates with the Go backend
 * via Wails runtime bindings. All Go methods are accessible through
 * window.go.main.App.*
 */

// Global state
let updateInfo = null;
let scannedFiles = [];
let isCopying = false;

/**
 * Initialize the application when DOM is ready.
 * Sets up event listeners for backend events and loads initial data.
 */
document.addEventListener('DOMContentLoaded', function () {
    // Listen for progress updates from Go backend
    if (typeof window.runtime !== 'undefined') {
        // Copy progress events
        window.runtime.EventsOn('copy:progress', handleProgressEvent);
        window.runtime.EventsOn('copy:complete', handleCompleteEvent);
        window.runtime.EventsOn('copy:cancelled', handleCancelledEvent);

        // Update progress events
        window.runtime.EventsOn('update:progress', function (message) {
            showToast(message, 'info');
        });
    }

    // Load initial data
    loadVersion();
    loadConfig();
    checkForUpdates();
});

/**
 * Load and display the current application version.
 * The version is shown in the header for user reference.
 */
async function loadVersion() {
    try {
        const version = await window.go.main.App.GetCurrentVersion();
        document.getElementById('versionText').textContent = version;
    } catch (err) {
        console.error('Failed to get version:', err);
    }
}

/**
 * Load the saved configuration and populate the form.
 * This restores user preferences from the last session.
 */
async function loadConfig() {
    try {
        const config = await window.go.main.App.GetConfig();
        if (config) {
            document.getElementById('sourcePath').value = config.source || '';
            document.getElementById('destPath').value = config.destination || '';
            document.getElementById('workers').value = config.workers || 10;
            document.getElementById('extensions').value = (config.extensions || []).join(',');
            document.getElementById('dryRun').checked = config.dryRun || false;
        }
    } catch (err) {
        console.error('Failed to load config:', err);
    }
}

/**
 * Check GitHub for available updates.
 * Shows the update button with a pulse animation if an update is available.
 */
async function checkForUpdates() {
    try {
        updateInfo = await window.go.main.App.CheckForUpdate();

        if (updateInfo && updateInfo.available) {
            const updateBtn = document.getElementById('updateBtn');
            updateBtn.classList.add('visible');
            updateBtn.title = `Update to ${updateInfo.latestVersion} available! Click to install.`;
            console.log('Update available:', updateInfo.latestVersion);
        } else {
            const updateBtn = document.getElementById('updateBtn');
            updateBtn.classList.remove('visible');
        }
    } catch (err) {
        // Network errors are expected when offline - fail silently
        console.error('Failed to check for updates:', err);
    }
}

/**
 * Download and install the available update.
 * This will restart the application after installing.
 */
async function performUpdate() {
    if (!updateInfo || !updateInfo.downloadUrl) {
        showToast('No update information available', 'error');
        return;
    }

    const updateBtn = document.getElementById('updateBtn');
    updateBtn.disabled = true;

    showToast(`Downloading ${updateInfo.latestVersion}...`, 'info');

    try {
        await window.go.main.App.PerformUpdate(updateInfo.downloadUrl);
        showToast('Update installed! Restarting...', 'success');
    } catch (err) {
        showToast('Update failed: ' + err, 'error');
        updateBtn.disabled = false;
    }
}

/**
 * Open native folder picker for source directory.
 * Updates the config when a folder is selected.
 */
async function selectSource() {
    try {
        const path = await window.go.main.App.SelectSourceFolder();
        if (path) {
            document.getElementById('sourcePath').value = path;
            await updateConfigFromForm();
            // Clear previous scan results
            scannedFiles = [];
            updateFileCount();
            disableCopyButtons();
        }
    } catch (err) {
        showToast('Error selecting folder: ' + err, 'error');
    }
}

/**
 * Open native folder picker for destination directory.
 * Updates the config when a folder is selected.
 */
async function selectDest() {
    try {
        const path = await window.go.main.App.SelectDestFolder();
        if (path) {
            document.getElementById('destPath').value = path;
            await updateConfigFromForm();
        }
    } catch (err) {
        showToast('Error selecting folder: ' + err, 'error');
    }
}

/**
 * Update the backend config with current form values.
 * This is called whenever the user changes a setting.
 */
async function updateConfigFromForm() {
    const extensions = document.getElementById('extensions').value
        .split(',')
        .map(e => e.trim())
        .filter(e => e.length > 0);

    const config = {
        source: document.getElementById('sourcePath').value,
        destination: document.getElementById('destPath').value,
        workers: parseInt(document.getElementById('workers').value) || 10,
        extensions: extensions,
        dryRun: document.getElementById('dryRun').checked,
        maxRetries: 3,
        overwrite: false
    };

    try {
        await window.go.main.App.UpdateConfig(config);
    } catch (err) {
        showToast('Failed to update config: ' + err, 'error');
    }
}

/**
 * Scan the source directory for files matching the filter.
 * Enables the copy buttons after a successful scan.
 */
async function scanFiles() {
    const sourcePath = document.getElementById('sourcePath').value;
    if (!sourcePath) {
        showToast('Please select a source folder first', 'error');
        return;
    }

    await updateConfigFromForm();

    const scanBtn = document.getElementById('scanBtn');
    scanBtn.disabled = true;
    scanBtn.innerHTML = `
        <svg class="spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" opacity="0.3"/>
            <path d="M12 2C6.48 2 2 6.48 2 12"/>
        </svg>
        Scanning...
    `;

    try {
        scannedFiles = await window.go.main.App.ScanFiles();
        updateFileCount();

        if (scannedFiles.length > 0) {
            enableCopyButtons();
            showToast(`Found ${scannedFiles.length} file(s) ready to copy`, 'success');
        } else {
            showToast('No files found matching the filter', 'info');
            disableCopyButtons();
        }
    } catch (err) {
        showToast('Scan failed: ' + err, 'error');
        disableCopyButtons();
    } finally {
        scanBtn.disabled = false;
        scanBtn.innerHTML = `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <path d="M21 21l-4.35-4.35"/>
            </svg>
            Scan Files
        `;
    }
}

/**
 * Start the copy operation.
 * @param {boolean} overwrite - Whether to overwrite existing files
 */
async function startCopy(overwrite) {
    const destPath = document.getElementById('destPath').value;
    if (!destPath) {
        showToast('Please select a destination folder', 'error');
        return;
    }

    if (scannedFiles.length === 0) {
        showToast('Please scan files first', 'error');
        return;
    }

    isCopying = true;
    disableCopyButtons();
    showProgressCard();
    hideResultsCard();
    resetProgress();

    try {
        const result = await window.go.main.App.StartCopy(overwrite);
        // Result is handled by the complete event
    } catch (err) {
        showToast('Copy failed: ' + err, 'error');
        hideProgressCard();
        enableCopyButtons();
        isCopying = false;
    }
}

/**
 * Cancel an ongoing copy operation.
 * Remaining files will not be copied.
 */
async function cancelCopy() {
    try {
        await window.go.main.App.CancelCopy();
        showToast('Copy operation cancelled', 'info');
    } catch (err) {
        console.error('Failed to cancel:', err);
    }
}

/**
 * Handle progress events from the backend.
 * Updates the progress bar and current file display.
 */
function handleProgressEvent(data) {
    const progressFill = document.getElementById('progressFill');
    const progressText = document.getElementById('progressText');
    const progressCount = document.getElementById('progressCount');
    const currentFile = document.getElementById('currentFile');

    progressFill.style.width = data.percent + '%';
    progressText.textContent = Math.round(data.percent) + '%';
    progressCount.textContent = `${data.current}/${data.total}`;
    currentFile.textContent = data.fileName;
}

/**
 * Handle the copy complete event.
 * Displays the results card with statistics.
 */
function handleCompleteEvent(result) {
    isCopying = false;
    hideProgressCard();
    showResultsCard(result);
    enableCopyButtons();

    if (result.success) {
        showToast(result.message, 'success');
    } else {
        showToast(result.message, 'error');
    }
}

/**
 * Handle the copy cancelled event.
 */
function handleCancelledEvent() {
    isCopying = false;
    hideProgressCard();
    enableCopyButtons();
}

// ===== UI Helper Functions =====

function updateFileCount() {
    const fileCount = document.getElementById('fileCount');
    if (scannedFiles.length > 0) {
        fileCount.textContent = `📁 ${scannedFiles.length} file(s) found`;
    } else {
        fileCount.textContent = '';
    }
}

function enableCopyButtons() {
    document.getElementById('copyOverwriteBtn').disabled = false;
    document.getElementById('copySkipBtn').disabled = false;
}

function disableCopyButtons() {
    document.getElementById('copyOverwriteBtn').disabled = true;
    document.getElementById('copySkipBtn').disabled = true;
}

function showProgressCard() {
    document.getElementById('progressCard').classList.add('active');
}

function hideProgressCard() {
    document.getElementById('progressCard').classList.remove('active');
}

function showResultsCard(result) {
    const card = document.getElementById('resultsCard');
    card.classList.add('active');

    document.getElementById('resultSuccess').textContent = result.successful;
    document.getElementById('resultFailed').textContent = result.failed;
    document.getElementById('resultSkipped').textContent = result.skipped;
    document.getElementById('resultDuration').textContent = result.duration.toFixed(2) + 's';
}

function hideResultsCard() {
    document.getElementById('resultsCard').classList.remove('active');
}

function resetProgress() {
    document.getElementById('progressFill').style.width = '0%';
    document.getElementById('progressText').textContent = '0%';
    document.getElementById('progressCount').textContent = '0/0';
    document.getElementById('currentFile').textContent = '';
}

/**
 * Display a toast notification.
 * Toasts auto-dismiss after 4 seconds.
 * 
 * @param {string} message - The message to display
 * @param {string} type - 'success', 'error', or 'info'
 */
function showToast(message, type) {
    const container = document.getElementById('toast-container');

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    let icon = '✓';
    if (type === 'error') icon = '✗';
    if (type === 'info') icon = 'ℹ';

    toast.innerHTML = `
        <span class="toast-icon">${icon}</span>
        <span class="toast-message">${message}</span>
    `;

    container.appendChild(toast);

    // Auto-remove after 4 seconds
    setTimeout(() => {
        toast.classList.add('hiding');
        toast.addEventListener('transitionend', () => {
            if (toast.parentElement) {
                toast.remove();
            }
        });
    }, 4000);
}

// Add CSS for spin animation
const style = document.createElement('style');
style.textContent = `
    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
    .spin {
        animation: spin 1s linear infinite;
    }
`;
document.head.appendChild(style);
