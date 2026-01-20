export namespace config {
	
	export class Destination {
	    id: string;
	    path: string;
	    overwrite: boolean;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Destination(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.overwrite = source["overwrite"];
	        this.enabled = source["enabled"];
	    }
	}
	export class CopyGroup {
	    id: string;
	    name: string;
	    source: string;
	    destinations: Destination[];
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CopyGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.source = source["source"];
	        this.destinations = this.convertValues(source["destinations"], Destination);
	        this.enabled = source["enabled"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Config {
	    source: string;
	    destination: string;
	    groups: CopyGroup[];
	    workers: number;
	    overwrite: boolean;
	    extensions: string[];
	    maxRetries: number;
	    dryRun: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source = source["source"];
	        this.destination = source["destination"];
	        this.groups = this.convertValues(source["groups"], CopyGroup);
	        this.workers = source["workers"];
	        this.overwrite = source["overwrite"];
	        this.extensions = source["extensions"];
	        this.maxRetries = source["maxRetries"];
	        this.dryRun = source["dryRun"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace main {
	
	export class CopyResult {
	    success: boolean;
	    message: string;
	    totalFiles: number;
	    successful: number;
	    failed: number;
	    skipped: number;
	    failedFiles: string[];
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new CopyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.totalFiles = source["totalFiles"];
	        this.successful = source["successful"];
	        this.failed = source["failed"];
	        this.skipped = source["skipped"];
	        this.failedFiles = source["failedFiles"];
	        this.duration = source["duration"];
	    }
	}
	export class UpdateInfo {
	    available: boolean;
	    currentVersion: string;
	    latestVersion: string;
	    downloadUrl: string;
	    releaseUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.downloadUrl = source["downloadUrl"];
	        this.releaseUrl = source["releaseUrl"];
	    }
	}

}

