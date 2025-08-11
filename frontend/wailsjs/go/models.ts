export namespace common_model {
	
	export class CommonBytesResp {
	    status: boolean;
	    bytes: number[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new CommonBytesResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.bytes = source["bytes"];
	        this.err = source["err"];
	    }
	}
	export class CommonConfigOption {
	    title: string;
	    key: string;
	    web_key: string;
	    default_value: any;
	    default_value_dynamic?: boolean;
	    tips?: string;
	    default_options?: string[];
	    auto_setting?: boolean;
	    form_required?: boolean;
	    form_hidden: boolean;
	    form_component_type?: number;
	    form_config_type?: string;
	    show_in_table: boolean;
	    col_width?: number;
	    col_fixed?: string;
	    col_hidden?: boolean;
	    edit?: boolean;
	    sort?: string;
	    custom_filter_dropdown?: boolean;
	    row_group?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CommonConfigOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.key = source["key"];
	        this.web_key = source["web_key"];
	        this.default_value = source["default_value"];
	        this.default_value_dynamic = source["default_value_dynamic"];
	        this.tips = source["tips"];
	        this.default_options = source["default_options"];
	        this.auto_setting = source["auto_setting"];
	        this.form_required = source["form_required"];
	        this.form_hidden = source["form_hidden"];
	        this.form_component_type = source["form_component_type"];
	        this.form_config_type = source["form_config_type"];
	        this.show_in_table = source["show_in_table"];
	        this.col_width = source["col_width"];
	        this.col_fixed = source["col_fixed"];
	        this.col_hidden = source["col_hidden"];
	        this.edit = source["edit"];
	        this.sort = source["sort"];
	        this.custom_filter_dropdown = source["custom_filter_dropdown"];
	        this.row_group = source["row_group"];
	    }
	}
	export class CommonConfigOptionsResp {
	    status: boolean;
	    msg: string;
	    options: CommonConfigOption[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new CommonConfigOptionsResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.msg = source["msg"];
	        this.options = this.convertValues(source["options"], CommonConfigOption);
	        this.err = source["err"];
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
	export class CommonReq {
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new CommonReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.msg = source["msg"];
	    }
	}
	export class CommonResp {
	    status: boolean;
	    msg: string;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new CommonResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.msg = source["msg"];
	        this.err = source["err"];
	    }
	}
	export class DBInfo {
	    info: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new DBInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.info = source["info"];
	    }
	}
	export class DBManageReq {
	    module: string;
	    operation: string;
	    conditions: DBInfo[];
	    info: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new DBManageReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.module = source["module"];
	        this.operation = source["operation"];
	        this.conditions = this.convertValues(source["conditions"], DBInfo);
	        this.info = source["info"];
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
	export class DBManageResp {
	    status: boolean;
	    msg: number[];
	    err: string;
	    infos: DBInfo[];
	
	    static createFrom(source: any = {}) {
	        return new DBManageResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.msg = source["msg"];
	        this.err = source["err"];
	        this.infos = this.convertValues(source["infos"], DBInfo);
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
	export class DynamicFormReq {
	    change_key: string;
	    trick_key: string;
	
	    static createFrom(source: any = {}) {
	        return new DynamicFormReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.change_key = source["change_key"];
	        this.trick_key = source["trick_key"];
	    }
	}
	export class DynamicFormResp {
	    status: boolean;
	    err: string;
	    values: string[];
	
	    static createFrom(source: any = {}) {
	        return new DynamicFormResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.err = source["err"];
	        this.values = source["values"];
	    }
	}
	export class LanguageEncodeReq {
	    data: string;
	    src_charset: string;
	    dst_charset: string;
	
	    static createFrom(source: any = {}) {
	        return new LanguageEncodeReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.src_charset = source["src_charset"];
	        this.dst_charset = source["dst_charset"];
	    }
	}
	export class LanguageEncodeResp {
	    status: boolean;
	    msg: string;
	    err: string;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new LanguageEncodeResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.msg = source["msg"];
	        this.err = source["err"];
	        this.data = source["data"];
	    }
	}
	export class TaskPauseReq {
	    task_id: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskPauseReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.task_id = source["task_id"];
	        this.status = source["status"];
	    }
	}
	export class TaskStatus {
	    task_type: string;
	    status: string;
	    task_id: string;
	    params: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new TaskStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.task_type = source["task_type"];
	        this.status = source["status"];
	        this.task_id = source["task_id"];
	        this.params = source["params"];
	    }
	}
	export class TaskStatusResp {
	    status: boolean;
	    task_status: TaskStatus;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskStatusResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.task_status = this.convertValues(source["task_status"], TaskStatus);
	        this.err = source["err"];
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
	export class TasksResp {
	    status: boolean;
	    tasks: TaskStatus[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new TasksResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.tasks = this.convertValues(source["tasks"], TaskStatus);
	        this.err = source["err"];
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
	
	export class AppStatus {
	    run_mode: string;
	    memory_usage: string;
	    cpu_percent: string;
	
	    static createFrom(source: any = {}) {
	        return new AppStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.run_mode = source["run_mode"];
	        this.memory_usage = source["memory_usage"];
	        this.cpu_percent = source["cpu_percent"];
	    }
	}

}

export namespace models {
	
	export class AgentAttribute {
	    hp?: number;
	    attack?: number;
	    defense?: number;
	    impact?: number;
	    critical_rate?: number;
	    critical_damage?: number;
	    anomaly_mastery?: number;
	    anomaly_proficiency?: number;
	    penetration_radio?: number;
	    penetration?: number;
	    energy_regen?: number;
	    sheer_force?: number;
	    hp_bonus?: number;
	    defense_bonus?: number;
	    attack_bonus?: number;
	    common_damage_bonus?: number;
	    ice_damage_bonus?: number;
	    electric_damage_bonus?: number;
	    physical_damage_bonus?: number;
	    fire_damage_bonus?: number;
	    ether_damage_bonus?: number;
	    common_sheer_damage_bonus?: number;
	    ice_sheer_damage_bonus?: number;
	    electric_sheer_damage_bonus?: number;
	    physical_sheer_damage_bonus?: number;
	    fire_sheer_damage_bonus?: number;
	    ether_sheer_damage_bonus?: number;
	    defense_reduction?: number;
	    damage_resistance?: number;
	    stun_damage_multiplier?: number;
	
	    static createFrom(source: any = {}) {
	        return new AgentAttribute(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hp = source["hp"];
	        this.attack = source["attack"];
	        this.defense = source["defense"];
	        this.impact = source["impact"];
	        this.critical_rate = source["critical_rate"];
	        this.critical_damage = source["critical_damage"];
	        this.anomaly_mastery = source["anomaly_mastery"];
	        this.anomaly_proficiency = source["anomaly_proficiency"];
	        this.penetration_radio = source["penetration_radio"];
	        this.penetration = source["penetration"];
	        this.energy_regen = source["energy_regen"];
	        this.sheer_force = source["sheer_force"];
	        this.hp_bonus = source["hp_bonus"];
	        this.defense_bonus = source["defense_bonus"];
	        this.attack_bonus = source["attack_bonus"];
	        this.common_damage_bonus = source["common_damage_bonus"];
	        this.ice_damage_bonus = source["ice_damage_bonus"];
	        this.electric_damage_bonus = source["electric_damage_bonus"];
	        this.physical_damage_bonus = source["physical_damage_bonus"];
	        this.fire_damage_bonus = source["fire_damage_bonus"];
	        this.ether_damage_bonus = source["ether_damage_bonus"];
	        this.common_sheer_damage_bonus = source["common_sheer_damage_bonus"];
	        this.ice_sheer_damage_bonus = source["ice_sheer_damage_bonus"];
	        this.electric_sheer_damage_bonus = source["electric_sheer_damage_bonus"];
	        this.physical_sheer_damage_bonus = source["physical_sheer_damage_bonus"];
	        this.fire_sheer_damage_bonus = source["fire_sheer_damage_bonus"];
	        this.ether_sheer_damage_bonus = source["ether_sheer_damage_bonus"];
	        this.defense_reduction = source["defense_reduction"];
	        this.damage_resistance = source["damage_resistance"];
	        this.stun_damage_multiplier = source["stun_damage_multiplier"];
	    }
	}
	export class AgentFeatures {
	    LifeDestroy: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AgentFeatures(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LifeDestroy = source["LifeDestroy"];
	    }
	}
	export class CommonAttribute {
	    out_game: AgentAttribute;
	    in_game: AgentAttribute;
	    buff: AgentAttribute;
	
	    static createFrom(source: any = {}) {
	        return new CommonAttribute(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.out_game = this.convertValues(source["out_game"], AgentAttribute);
	        this.in_game = this.convertValues(source["in_game"], AgentAttribute);
	        this.buff = this.convertValues(source["buff"], AgentAttribute);
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
	export class TestData {
	    level_base: number;
	    monster_base_defense: number;
	    damage_multiplier: number;
	    attribute: CommonAttribute;
	
	    static createFrom(source: any = {}) {
	        return new TestData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level_base = source["level_base"];
	        this.monster_base_defense = source["monster_base_defense"];
	        this.damage_multiplier = source["damage_multiplier"];
	        this.attribute = this.convertValues(source["attribute"], CommonAttribute);
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
	export class DamageFuzzParam {
	    name: string;
	    star: number;
	    engine: string;
	    engine_star: number;
	    improve: boolean;
	    stun: boolean;
	    test_data: TestData;
	    driver_type: string[];
	    driver_ids: string[];
	
	    static createFrom(source: any = {}) {
	        return new DamageFuzzParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.star = source["star"];
	        this.engine = source["engine"];
	        this.engine_star = source["engine_star"];
	        this.improve = source["improve"];
	        this.stun = source["stun"];
	        this.test_data = this.convertValues(source["test_data"], TestData);
	        this.driver_type = source["driver_type"];
	        this.driver_ids = source["driver_ids"];
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
	export class DriverDiskMainStat {
	    HP: number;
	    HPBonus: number;
	    Attack: number;
	    Defense: number;
	    DefenseBonus: number;
	    CriticalDamage: number;
	    CriticalRate: number;
	    AnomalyProficiency: number;
	    CommonDamageBonus: number;
	    IceDamageBonus: number;
	    ElectricDamageBonus: number;
	    PhysicalDamageBonus: number;
	    FireDamageBonus: number;
	    EtherDamageBonus: number;
	    CommonSheerDamageBonus: number;
	    IceSheerDamageBonus: number;
	    ElectricSheerDamageBonus: number;
	    PhysicalSheerDamageBonus: number;
	    FireSheerDamageBonus: number;
	    EtherSheerDamageBonus: number;
	    AttackBonus: number;
	    PenetrationRadio: number;
	    AnomalyMastery: number;
	    EnergyRegen: number;
	
	    static createFrom(source: any = {}) {
	        return new DriverDiskMainStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HP = source["HP"];
	        this.HPBonus = source["HPBonus"];
	        this.Attack = source["Attack"];
	        this.Defense = source["Defense"];
	        this.DefenseBonus = source["DefenseBonus"];
	        this.CriticalDamage = source["CriticalDamage"];
	        this.CriticalRate = source["CriticalRate"];
	        this.AnomalyProficiency = source["AnomalyProficiency"];
	        this.CommonDamageBonus = source["CommonDamageBonus"];
	        this.IceDamageBonus = source["IceDamageBonus"];
	        this.ElectricDamageBonus = source["ElectricDamageBonus"];
	        this.PhysicalDamageBonus = source["PhysicalDamageBonus"];
	        this.FireDamageBonus = source["FireDamageBonus"];
	        this.EtherDamageBonus = source["EtherDamageBonus"];
	        this.CommonSheerDamageBonus = source["CommonSheerDamageBonus"];
	        this.IceSheerDamageBonus = source["IceSheerDamageBonus"];
	        this.ElectricSheerDamageBonus = source["ElectricSheerDamageBonus"];
	        this.PhysicalSheerDamageBonus = source["PhysicalSheerDamageBonus"];
	        this.FireSheerDamageBonus = source["FireSheerDamageBonus"];
	        this.EtherSheerDamageBonus = source["EtherSheerDamageBonus"];
	        this.AttackBonus = source["AttackBonus"];
	        this.PenetrationRadio = source["PenetrationRadio"];
	        this.AnomalyMastery = source["AnomalyMastery"];
	        this.EnergyRegen = source["EnergyRegen"];
	    }
	}
	export class DriverDiskSubStat {
	    Attack: number;
	    AttackBonus: number;
	    Penetration: number;
	    CriticalDamage: number;
	    CriticalRate: number;
	    AnomalyProficiency: number;
	    Defense: number;
	    DefenseBonus: number;
	    HP: number;
	    HPBonus: number;
	
	    static createFrom(source: any = {}) {
	        return new DriverDiskSubStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Attack = source["Attack"];
	        this.AttackBonus = source["AttackBonus"];
	        this.Penetration = source["Penetration"];
	        this.CriticalDamage = source["CriticalDamage"];
	        this.CriticalRate = source["CriticalRate"];
	        this.AnomalyProficiency = source["AnomalyProficiency"];
	        this.Defense = source["Defense"];
	        this.DefenseBonus = source["DefenseBonus"];
	        this.HP = source["HP"];
	        this.HPBonus = source["HPBonus"];
	    }
	}
	export class DriverDiskStat {
	    Name: string;
	    Position: number;
	    Main: DriverDiskMainStat;
	    Sub: DriverDiskSubStat;
	
	    static createFrom(source: any = {}) {
	        return new DriverDiskStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Position = source["Position"];
	        this.Main = this.convertValues(source["Main"], DriverDiskMainStat);
	        this.Sub = this.convertValues(source["Sub"], DriverDiskSubStat);
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

export namespace zzz_models {
	
	export class DriverFuzzResp {
	    err: string;
	    status: boolean;
	    out_game: string;
	    in_game: string;
	    disk1: string;
	    disk2: string;
	    disk3: string;
	    disk4: string;
	    disk5: string;
	    disk6: string;
	
	    static createFrom(source: any = {}) {
	        return new DriverFuzzResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.err = source["err"];
	        this.status = source["status"];
	        this.out_game = source["out_game"];
	        this.in_game = source["in_game"];
	        this.disk1 = source["disk1"];
	        this.disk2 = source["disk2"];
	        this.disk3 = source["disk3"];
	        this.disk4 = source["disk4"];
	        this.disk5 = source["disk5"];
	        this.disk6 = source["disk6"];
	    }
	}
	export class DriverParserResp {
	    status: boolean;
	    ids: string[];
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new DriverParserResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.ids = source["ids"];
	        this.err = source["err"];
	    }
	}
	export class FileInfo {
	    id: string;
	    data: number[];
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.data = source["data"];
	    }
	}
	export class TestProxyInfo {
	    name: string;
	    star: number;
	    engine: string;
	    engine_star: number;
	    driver_set: string;
	
	    static createFrom(source: any = {}) {
	        return new TestProxyInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.star = source["star"];
	        this.engine = source["engine"];
	        this.engine_star = source["engine_star"];
	        this.driver_set = source["driver_set"];
	    }
	}
	export class TestProxyBuffReq {
	    proxy1: TestProxyInfo;
	    proxy2: TestProxyInfo;
	
	    static createFrom(source: any = {}) {
	        return new TestProxyBuffReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.proxy1 = this.convertValues(source["proxy1"], TestProxyInfo);
	        this.proxy2 = this.convertValues(source["proxy2"], TestProxyInfo);
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
	export class TestProxyBuffResp {
	    err: string;
	    status: boolean;
	    attribute: models.AgentAttribute;
	
	    static createFrom(source: any = {}) {
	        return new TestProxyBuffResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.err = source["err"];
	        this.status = source["status"];
	        this.attribute = this.convertValues(source["attribute"], models.AgentAttribute);
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

