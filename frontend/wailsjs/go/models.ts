export namespace backend {
	
	export class BaseCalculationResult {
	    decimal: string;
	    hex: string;
	    binary: string;
	    groupedBinary: string;
	
	    static createFrom(source: any = {}) {
	        return new BaseCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.decimal = source["decimal"];
	        this.hex = source["hex"];
	        this.binary = source["binary"];
	        this.groupedBinary = source["groupedBinary"];
	    }
	}
	export class BaseConversionResult {
	    decimal: string;
	    hex: string;
	    binary: string;
	    groupedBinary: string;
	
	    static createFrom(source: any = {}) {
	        return new BaseConversionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.decimal = source["decimal"];
	        this.hex = source["hex"];
	        this.binary = source["binary"];
	        this.groupedBinary = source["groupedBinary"];
	    }
	}
	export class CIDRCalculationResult {
	    networkAddress: string;
	    broadcastAddress: string;
	    subnetMask: string;
	    wildcardMask: string;
	    firstHostAddress: string;
	    lastHostAddress: string;
	    totalHosts: number;
	    usableHosts: number;
	    cidr: string;
	    binarySubnetMask: string;
	    ipClass: string;
	    ipType: string;
	    inputIp: string;
	    inputWasHost: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CIDRCalculationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.networkAddress = source["networkAddress"];
	        this.broadcastAddress = source["broadcastAddress"];
	        this.subnetMask = source["subnetMask"];
	        this.wildcardMask = source["wildcardMask"];
	        this.firstHostAddress = source["firstHostAddress"];
	        this.lastHostAddress = source["lastHostAddress"];
	        this.totalHosts = source["totalHosts"];
	        this.usableHosts = source["usableHosts"];
	        this.cidr = source["cidr"];
	        this.binarySubnetMask = source["binarySubnetMask"];
	        this.ipClass = source["ipClass"];
	        this.ipType = source["ipType"];
	        this.inputIp = source["inputIp"];
	        this.inputWasHost = source["inputWasHost"];
	    }
	}
	export class CRLUpdateResult {
	    crlPem: string;
	    addedCount: number;
	    totalRevokedCount: number;
	    revokedSerialNumbers: string[];
	
	    static createFrom(source: any = {}) {
	        return new CRLUpdateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.crlPem = source["crlPem"];
	        this.addedCount = source["addedCount"];
	        this.totalRevokedCount = source["totalRevokedCount"];
	        this.revokedSerialNumbers = source["revokedSerialNumbers"];
	    }
	}
	export class ClientCertificate {
	    commonName: string;
	    certificatePem: string;
	    privateKeyPem: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientCertificate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.commonName = source["commonName"];
	        this.certificatePem = source["certificatePem"];
	        this.privateKeyPem = source["privateKeyPem"];
	    }
	}
	export class ConversionResult {
	    B: number;
	    KB: number;
	    MB: number;
	    GB: number;
	    TB: number;
	    KiB: number;
	    MiB: number;
	    GiB: number;
	    TiB: number;
	    bits: number;
	    Kbits: number;
	    Mbits: number;
	    Gbits: number;
	    Tbits: number;
	    Kibits: number;
	    Mibits: number;
	    Gibits: number;
	    Tibits: number;
	
	    static createFrom(source: any = {}) {
	        return new ConversionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.B = source["B"];
	        this.KB = source["KB"];
	        this.MB = source["MB"];
	        this.GB = source["GB"];
	        this.TB = source["TB"];
	        this.KiB = source["KiB"];
	        this.MiB = source["MiB"];
	        this.GiB = source["GiB"];
	        this.TiB = source["TiB"];
	        this.bits = source["bits"];
	        this.Kbits = source["Kbits"];
	        this.Mbits = source["Mbits"];
	        this.Gbits = source["Gbits"];
	        this.Tbits = source["Tbits"];
	        this.Kibits = source["Kibits"];
	        this.Mibits = source["Mibits"];
	        this.Gibits = source["Gibits"];
	        this.Tibits = source["Tibits"];
	    }
	}
	export class DataTransferRateResult {
	    "B/s": number;
	    "KB/s": number;
	    "MB/s": number;
	    "GB/s": number;
	    "TB/s": number;
	    "KiB/s": number;
	    "MiB/s": number;
	    "GiB/s": number;
	    "TiB/s": number;
	    "bit/s": number;
	    "Kbit/s": number;
	    "Mbit/s": number;
	    "Gbit/s": number;
	    "Tbit/s": number;
	    "Kibit/s": number;
	    "Mibit/s": number;
	    "Gibit/s": number;
	    "Tibit/s": number;
	
	    static createFrom(source: any = {}) {
	        return new DataTransferRateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this["B/s"] = source["B/s"];
	        this["KB/s"] = source["KB/s"];
	        this["MB/s"] = source["MB/s"];
	        this["GB/s"] = source["GB/s"];
	        this["TB/s"] = source["TB/s"];
	        this["KiB/s"] = source["KiB/s"];
	        this["MiB/s"] = source["MiB/s"];
	        this["GiB/s"] = source["GiB/s"];
	        this["TiB/s"] = source["TiB/s"];
	        this["bit/s"] = source["bit/s"];
	        this["Kbit/s"] = source["Kbit/s"];
	        this["Mbit/s"] = source["Mbit/s"];
	        this["Gbit/s"] = source["Gbit/s"];
	        this["Tbit/s"] = source["Tbit/s"];
	        this["Kibit/s"] = source["Kibit/s"];
	        this["Mibit/s"] = source["Mibit/s"];
	        this["Gibit/s"] = source["Gibit/s"];
	        this["Tibit/s"] = source["Tibit/s"];
	    }
	}
	export class JWTDecodeResult {
	    header: string;
	    payload: string;
	    valid: boolean;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new JWTDecodeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.header = source["header"];
	        this.payload = source["payload"];
	        this.valid = source["valid"];
	        this.error = source["error"];
	    }
	}
	export class MTLSCertificateResult {
	    caCertificatePem: string;
	    caPrivateKeyPem: string;
	    serverCertificatePem: string;
	    serverPrivateKeyPem: string;
	    clientCertificatePem: string;
	    clientPrivateKeyPem: string;
	
	    static createFrom(source: any = {}) {
	        return new MTLSCertificateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.caCertificatePem = source["caCertificatePem"];
	        this.caPrivateKeyPem = source["caPrivateKeyPem"];
	        this.serverCertificatePem = source["serverCertificatePem"];
	        this.serverPrivateKeyPem = source["serverPrivateKeyPem"];
	        this.clientCertificatePem = source["clientCertificatePem"];
	        this.clientPrivateKeyPem = source["clientPrivateKeyPem"];
	    }
	}
	export class MTLSCertificatesMultiClientResult {
	    caCertificatePem: string;
	    caPrivateKeyPem: string;
	    crlPem: string;
	    serverCertificatePem: string;
	    serverPrivateKeyPem: string;
	    clientCertificates: ClientCertificate[];
	
	    static createFrom(source: any = {}) {
	        return new MTLSCertificatesMultiClientResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.caCertificatePem = source["caCertificatePem"];
	        this.caPrivateKeyPem = source["caPrivateKeyPem"];
	        this.crlPem = source["crlPem"];
	        this.serverCertificatePem = source["serverCertificatePem"];
	        this.serverPrivateKeyPem = source["serverPrivateKeyPem"];
	        this.clientCertificates = this.convertValues(source["clientCertificates"], ClientCertificate);
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
	export class SelfSignedCertificateResult {
	    certificatePem: string;
	    privateKeyPem: string;
	
	    static createFrom(source: any = {}) {
	        return new SelfSignedCertificateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.certificatePem = source["certificatePem"];
	        this.privateKeyPem = source["privateKeyPem"];
	    }
	}
	export class TransferTimeResult {
	    seconds: number;
	    minutes: number;
	    hours: number;
	    days: number;
	
	    static createFrom(source: any = {}) {
	        return new TransferTimeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.seconds = source["seconds"];
	        this.minutes = source["minutes"];
	        this.hours = source["hours"];
	        this.days = source["days"];
	    }
	}

}

