export namespace disk {
	
	export class CloneOptions {
	    source: string;
	    destination: string;
	    bufferSize: number;
	
	    static createFrom(source: any = {}) {
	        return new CloneOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source = source["source"];
	        this.destination = source["destination"];
	        this.bufferSize = source["bufferSize"];
	    }
	}
	export class Info {
	    name: string;
	    path: string;
	    size: number;
	    type: string;
	    mountpoint: string;
	    fstype: string;
	    model: string;
	    serial: string;
	    health: string;
	    partitions?: Info[];
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.type = source["type"];
	        this.mountpoint = source["mountpoint"];
	        this.fstype = source["fstype"];
	        this.model = source["model"];
	        this.serial = source["serial"];
	        this.health = source["health"];
	        this.partitions = this.convertValues(source["partitions"], Info);
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
	
	export class SnapStatus {
	    isSnap: boolean;
	    hasBlockAccess: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SnapStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isSnap = source["isSnap"];
	        this.hasBlockAccess = source["hasBlockAccess"];
	    }
	}

}

