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
	    resolvedIcon: string;
	
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
	        this.resolvedIcon = source["resolvedIcon"];
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
	export class ConvertResult {
	    success: boolean;
	    query: string;
	    fromValue: string;
	    fromUnit: string;
	    toValue: string;
	    toUnit: string;
	    result: string;
	    category: string;
	    formatted: string;
	
	    static createFrom(source: any = {}) {
	        return new ConvertResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.query = source["query"];
	        this.fromValue = source["fromValue"];
	        this.fromUnit = source["fromUnit"];
	        this.toValue = source["toValue"];
	        this.toUnit = source["toUnit"];
	        this.result = source["result"];
	        this.category = source["category"];
	        this.formatted = source["formatted"];
	    }
	}
	export class EmojiItem {
	    emoji: string;
	    name: string;
	    shortcode: string;
	    category: string;
	    keywords: string[];
	
	    static createFrom(source: any = {}) {
	        return new EmojiItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.emoji = source["emoji"];
	        this.name = source["name"];
	        this.shortcode = source["shortcode"];
	        this.category = source["category"];
	        this.keywords = source["keywords"];
	    }
	}
	export class GuideShortcut {
	    key: string;
	    desc: string;
	
	    static createFrom(source: any = {}) {
	        return new GuideShortcut(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.desc = source["desc"];
	    }
	}
	export class GuideTopic {
	    id: string;
	    title: string;
	    category: string;
	    summary: string;
	    content: string[];
	    shortcuts: GuideShortcut[];
	    examples: string[];
	
	    static createFrom(source: any = {}) {
	        return new GuideTopic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.category = source["category"];
	        this.summary = source["summary"];
	        this.content = source["content"];
	        this.shortcuts = this.convertValues(source["shortcuts"], GuideShortcut);
	        this.examples = source["examples"];
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
	export class ScriptCommand {
	    id: string;
	    title: string;
	    path: string;
	    mode: string;
	    description: string;
	    author: string;
	
	    static createFrom(source: any = {}) {
	        return new ScriptCommand(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.path = source["path"];
	        this.mode = source["mode"];
	        this.description = source["description"];
	        this.author = source["author"];
	    }
	}
	export class TodoItem {
	    id: string;
	    title: string;
	    completed: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    completedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new TodoItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.completed = source["completed"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
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

