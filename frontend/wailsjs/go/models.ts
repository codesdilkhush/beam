export namespace backend {
	
	export class AppInfo {
	    id: string;
	    name: string;
	    comment: string;
	    exec: string;
	    icon: string;
	    iconPath: string;
	    terminal: boolean;
	    desktopPath: string;
	    score: number;
	    launchCount: number;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.comment = source["comment"];
	        this.exec = source["exec"];
	        this.icon = source["icon"];
	        this.iconPath = source["iconPath"];
	        this.terminal = source["terminal"];
	        this.desktopPath = source["desktopPath"];
	        this.score = source["score"];
	        this.launchCount = source["launchCount"];
	    }
	}
	export class CalcResult {
	    expression: string;
	    result: string;
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CalcResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.expression = source["expression"];
	        this.result = source["result"];
	        this.success = source["success"];
	    }
	}
	export class ClipboardItem {
	    id: string;
	    content: string;
	    preview: string;
	    // Go type: time
	    timestamp: any;
	    charCount: number;
	
	    static createFrom(source: any = {}) {
	        return new ClipboardItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	        this.preview = source["preview"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.charCount = source["charCount"];
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

