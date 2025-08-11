import React, {useEffect, useMemo, useRef, useState} from "react";
import PageWrapTabs from "@/pages/components/PageWrapTabs";
import {Allotment} from "allotment";
import {ProForm, ProFormDigit, ProFormList} from "@ant-design/pro-form";
import {Button, Drawer, Flex, Form, Image, InputNumber, message, Select} from "antd";
import globalStyles from "@/utils/style.module.less";
import {ProFormCheckbox, ProFormSelect} from "@ant-design/pro-components";
import L6CodeEditor from "@/components/level6/L6CodeEditor";
import L6Row from "@/components/level6/L6Layout/L6Row";
import L6Col from "@/components/level6/L6Layout/L6Col";
import { Collapse } from "antd/lib";
import {models, zzz_models} from "../../../../wailsjs/go/models";
import CommonAttribute = models.CommonAttribute;
import DamageFuzzParam = models.DamageFuzzParam;
import AgentAttribute = models.AgentAttribute;
import {DriverFuzz, GetProxyBuff, ReadDriverCache} from "../../../../wailsjs/go/main_control/Control";
import DriverFuzzResp = zzz_models.DriverFuzzResp;
import {EventsOff, EventsOn} from "../../../../wailsjs/runtime";
import {uuid} from "@/utils/uuid";
import DriverSelectPage, {DriverSelectRef} from "@/pages/zzz/driver_fuzz/driver_select";
import DriverHistoryPage from "@/pages/zzz/driver_fuzz/driver_history";
import type {UpdateDataType} from "@/pages/components/DBTable/types";
import DriverLogPage, {DriverLogRef} from "@/pages/zzz/driver_fuzz/log";


const formInitValues = {
    star: 0,
    engine_star: 1,
    test_data: {
        level_base: 794,
        monster_base_defense: 60,
        damage_multiplier: 1,
    },

    tmp_proxy1: "",
    tmp_proxy1_star: 0,
    tmp_proxy1_engine: "",
    tmp_proxy1_engine_star: 1,

    tmp_proxy2: "",
    tmp_proxy2_star: 0,
    tmp_proxy2_engine: "",
    tmp_proxy2_engine_star: 1,
}
const proxys = ["", "艾莲·乔", "星见雅", "仪玄", "苍角", "耀嘉音"]
const engines = ["","深海访客",  "霰落星殿", "青溟笼舍", "好斗的阿炮"]

const drivers = ["",'折枝剑歌', '河豚电音', '啄木鸟电音', '极地重金属', "静听佳音", "山大王", "摇摆爵士", "原始朋克"]







const DriverLogPageMemo = React.memo(DriverLogPage);
const DriverSelectPageMemo = React.memo(DriverSelectPage);



const Content: React.FC<{id?: string}> = ({id}) => {
    const [form] = Form.useForm();
    const driverSelectRef = useRef<DriverSelectRef>(null);
    const driverLogRef = useRef<DriverLogRef>(null);


    const eventID = useMemo(() => uuid(),[])
    const [showMore, setShowMore] = useState(false);
    const [fuzzResult, setFuzzResult] = useState<DriverFuzzResp>(new DriverFuzzResp())
    const [driverOpen, setDriverOpen] = useState(false)
    const [driverHistoryOpen, setDriverHistoryOpen] = useState(false)


    const [board, setBoard] = useState("")



    const testProxyDataHandle = async () => {

        const resp = await GetProxyBuff(new zzz_models.TestProxyBuffReq({
            proxy1: {
                name: form.getFieldValue("tmp_proxy1"),
                star: form.getFieldValue("tmp_proxy1_star"),
                engine: form.getFieldValue("tmp_proxy1_engine"),
                engine_star: form.getFieldValue("tmp_proxy1_engine_star"),
                driver_set: form.getFieldValue("tmp_proxy1_driver"),
            },
            proxy2: {
                name: form.getFieldValue("tmp_proxy2"),
                star: form.getFieldValue("tmp_proxy2_star"),
                engine: form.getFieldValue("tmp_proxy2_engine"),
                engine_star: form.getFieldValue("tmp_proxy2_engine_star"),
                driver_set: form.getFieldValue("tmp_proxy2_driver"),
            },
        }))
        if (!resp.status) {
            message.error(resp.err)
            return
        }
        const tmpData = []
        const attr = resp.attribute as Record<string, any>
        for (const key in attr) {
            const value = attr[key]
            if (value) {
                tmpData.push({
                    game_state: 'in_game',
                    attribute: key,
                    value: value
                })
            }
        }
        form.setFieldValue("tmp_data", tmpData)
        message.success("生成成功")
    }

    const driverHistoryDoubleClickHandle = async (currentRow: Record<string, any>, update: UpdateDataType) => {
        const fuzzResultResp = new DriverFuzzResp()

        {
            const resp = await  ReadDriverCache(`${currentRow.disk1}`)
            if (!resp.status) {
                message.error(`disk1: ${resp.err}`)
                return
            }
            fuzzResultResp.disk1 = resp.bytes.toString()
        }
        {
            const resp = await  ReadDriverCache(`${currentRow.disk2}`)
            if (!resp.status) {
                message.error(`disk2: ${resp.err}`)
                return
            }
            fuzzResultResp.disk2 = resp.bytes.toString()
        }
        {
            const resp = await  ReadDriverCache(`${currentRow.disk3}`)
            if (!resp.status) {
                message.error(`disk3: ${resp.err}`)
                return
            }
            fuzzResultResp.disk3 = resp.bytes.toString()
        }
        {
            const resp = await  ReadDriverCache(`${currentRow.disk4}`)
            if (!resp.status) {
                message.error(`disk4: ${resp.err}`)
                return
            }
            fuzzResultResp.disk4 = resp.bytes.toString()
        }
        {
            const resp = await  ReadDriverCache(`${currentRow.disk5}`)
            if (!resp.status) {
                message.error(`disk5: ${resp.err}`)
                return
            }
            fuzzResultResp.disk5 = resp.bytes.toString()
        }
        {
            const resp = await  ReadDriverCache(`${currentRow.disk6}`)
            if (!resp.status) {
                message.error(`disk6: ${resp.err}`)
                return
            }
            fuzzResultResp.disk6 = resp.bytes.toString()
        }

        fuzzResultResp.out_game = currentRow.fuzz_result
        setBoard(`${currentRow.fuzz_result}`)
        driverLogRef.current?.setLog(currentRow.fuzz_param)

        setFuzzResult(fuzzResultResp)
        setDriverHistoryOpen(false)
        message.success("获取历史成功")

    }

    return (<Allotment className={globalStyles['l6-frame']}>
        <Allotment.Pane minSize={100} preferredSize={`36%`}>
            <ProForm
                style={{height: 'calc( 100% - 36px )'}}
                form={form}
                layout={'horizontal'}
                labelAlign="left"
                labelCol={{span: 6}}
                autoFocusFirstInput
                initialValues={formInitValues}
                onFinish={async function (values) {
                    let driverIds = []
                    if (driverSelectRef.current) {
                        const rows = driverSelectRef.current?.tableRef?.api.getSelectedRows()
                        if (rows) {
                            driverIds = rows.map((item) => item.id)
                        }
                    }
                    if (values.tmp_data) {
                        const commonAttribute = new CommonAttribute()
                        commonAttribute.in_game = new AgentAttribute()
                        commonAttribute.out_game = new AgentAttribute()
                        for (let data of (values.tmp_data as any[])) {
                            if (data.game_state === 'out_game') {
                                // @ts-ignore
                                if (commonAttribute.out_game[data.attribute]) {
                                    // @ts-ignore
                                    commonAttribute.out_game[data.attribute] += parseFloat(data.value)
                                } else {
                                    // @ts-ignore
                                    commonAttribute.out_game[data.attribute] = parseFloat(data.value)
                                }
                            } else {
                                // @ts-ignore
                                if (commonAttribute.in_game[data.attribute]) {
                                    // @ts-ignore
                                    commonAttribute.in_game[data.attribute] += parseFloat(data.value)
                                } else {
                                    // @ts-ignore
                                    commonAttribute.in_game[data.attribute] = parseFloat(data.value)
                                }
                            }
                        }
                        values.test_data.attribute = commonAttribute
                    }
                    values.driver_ids = driverIds
                    const resp = await DriverFuzz(eventID, new DamageFuzzParam(values))
                    if (!resp.status) {
                        message.error(resp.err)
                        return
                    }
                    setFuzzResult(resp)
                    if (resp.out_game && !resp.in_game) {
                        setBoard(`${resp.out_game}`)
                    } else {
                        setBoard(`局外:\n${resp.out_game}\n局内:\n${resp.in_game}`)
                    }
                    message.success("计算结束")

                }}
                submitter={false}
            >
                <div className={globalStyles['l6-split']} style={{overflow: 'auto'}}>
                    <Flex gap={"small"}>
                        <ProForm.Item
                            name={"name"}
                            label={"代理人"}
                            initialValue={""}
                            labelCol={{span: 7}}
                            rules={[{required: true}]}
                            style={{width: '100%'}}
                        >
                            <Select
                                options={proxys.map((item) => ({value: item}))}
                            />
                        </ProForm.Item>
                        <ProForm.Item
                            name={"star"}
                            rules={[{required: true}]}
                        >
                            <Select options={
                                [0, 1,2,3,4,5,6]
                                    .map((item) => ({value: item}))}/>
                        </ProForm.Item>
                    </Flex>

                    <Flex gap={"small"} >
                        <ProForm.Item
                            name={"engine"}
                            label={"音擎"}
                            labelCol={{span: 7}}
                            rules={[{required: true}]}
                            style={{width: '100%'}}
                        >
                            <Select
                                options={engines.map((item) => ({value: item}))}
                            />
                        </ProForm.Item>
                        <ProForm.Item
                            name={"engine_star"}
                            rules={[{required: true}]}
                        >
                            <Select  options={
                                [1,2,3,4,5]
                                    .map((item) => ({value: item}))}/>
                        </ProForm.Item>
                    </Flex>

                    <ProFormSelect
                        name={"driver_type"}
                        label="驱动盘过滤"
                        valueEnum={Object.fromEntries(drivers.map(key => [key, key]))}
                        fieldProps={{
                            mode: 'multiple',
                        }}
                        rules={[{required: false}]}
                    />
                    <ProFormCheckbox
                        name={"stun"}
                        label={"失衡易伤"}
                    />


                    <Collapse items={[{
                        key: '1',
                        forceRender: true,
                        label: <Button ghost type={"primary"} onClick={()=> {
                            setShowMore(!showMore)
                        }}>{showMore?"隐藏测试属性":"显示测试属性"}</Button>,
                        children: <>
                            <ProFormSelect
                                name={["test_data", "level_base"]}
                                label={"等级基数"}
                                rules={[{required: true}]}
                                options={[
                                    {value: '794', label: '794'},
                                ]}
                            />
                            <ProFormSelect
                                name={["test_data", "monster_base_defense"]}
                                label={"怪物基础防御"}
                                rules={[{required: true}]}
                                options={[
                                    {value: '60', label: '60'},
                                ]}
                            />
                            <ProForm.Item
                                name={["test_data", "damage_multiplier"]}
                                label={"技能伤害倍率"}
                                rules={[{required: true}]}
                            >
                                <ProFormDigit />
                            </ProForm.Item>
                        <ProFormList
                            name={"tmp_data"}
                            className={globalStyles['l6-pro-form-list']}
                            creatorButtonProps={{
                                creatorButtonText: '添加参数'
                            }}
                            labelCol={{ span: 24 }}
                            wrapperCol={{ span: 24 }}
                        >
                            <div className={globalStyles['l6-pro-form-group']}>
                                <L6Row>
                                    <L6Col span={6}>
                                        <ProForm.Item
                                            name={"game_state"}
                                            label={"状态"}
                                            initialValue={"in_game"}
                                            labelCol={{ span: 24 }}
                                            wrapperCol={{ span: 24 }}
                                            rules={[{required: true}]}
                                        >
                                            <Select
                                                options={[
                                                    {value: 'out_game', label: '局外'},
                                                    {value: 'in_game', label: '局内'},
                                                ]}
                                            />
                                        </ProForm.Item>
                                    </L6Col>
                                    <L6Col span={12}>
                                        <ProForm.Item
                                            name={"attribute"}
                                            label={"属性"}
                                            initialValue={""}
                                            labelCol={{ span: 24 }}
                                            wrapperCol={{ span: 24 }}
                                            rules={[{required: true}]}
                                        >
                                            <Select
                                                options={[
                                                    {value: 'hp', label: '生命值'},
                                                    {value: 'attack', label: '攻击力'},
                                                    {value: 'defense', label: '防御力'},
                                                    {value: 'impact', label: '冲击力'},
                                                    {value: 'critical_rate', label: '暴击率'},
                                                    {value: 'critical_damage', label: '暴击伤害'},
                                                    {value: 'anomaly_mastery', label: '异常掌控'},
                                                    {value: 'anomaly_proficiency', label: '异常精通'},
                                                    {value: 'penetration_radio', label: '穿透率'},
                                                    {value: 'penetration', label: '穿透值'},
                                                    {value: 'energy_regen', label: '能量回复'},
                                                    {value: 'sheer_force', label: '贯穿力'},

                                                    {value: 'hp_bonus', label: '生命值加成'},
                                                    {value: 'defense_bonus', label: '防御力加成'},
                                                    {value: 'attack_bonus', label: '攻击力加成'},
                                                    {value: 'common_damage_bonus', label: '通用伤害加成'},
                                                    {value: 'ice_damage_bonus', label: '冰属性伤害加成'},
                                                    {value: 'electric_damage_bonus', label: '电属性伤害加成'},
                                                    {value: 'fire_damage_bonus', label: '火属性伤害加成'},
                                                    {value: 'ether_damage_bonus', label: '以太伤害加成'},

                                                    {value: 'common_sheer_damage_bonus', label: '贯穿伤害加成'},
                                                    {value: 'ice_sheer_damage_bonus', label: '贯穿冰属性伤害加成'},
                                                    {value: 'electric_sheer_damage_bonus', label: '贯穿电属性伤害加成'},
                                                    {value: 'physical_sheer_damage_bonus', label: '贯穿物理属性伤害加成'},
                                                    {value: 'fire_sheer_damage_bonus', label: '贯穿火属性伤害加成'},
                                                    {value: 'ether_sheer_damage_bonus', label: '贯穿以太伤害加成'},

                                                    {value: 'defense_reduction', label: '防御减伤'},
                                                    {value: 'damage_resistance', label: '抗性'},
                                                    {value: 'stun_damage_multiplier', label: '失衡易伤'},
                                                ]}
                                            />
                                        </ProForm.Item>
                                    </L6Col>

                                    <L6Col span={5}>
                                        <ProForm.Item
                                            name={"value"}
                                            label={"数值"}
                                            initialValue={""}
                                            labelCol={{ span: 24 }}
                                            wrapperCol={{ span: 24 }}
                                            rules={[{required: true}]}
                                        >
                                            <InputNumber   style={{ width: 70 }}/>
                                        </ProForm.Item>
                                    </L6Col>
                                </L6Row>

                            </div>
                        </ProFormList>
                        </>,
                    }]}  size={'small'} bordered={false} ghost />
                    <Collapse items={[{
                        key: '2',
                        forceRender: true,
                        label: <Button ghost type={"primary"} onClick={()=> {
                        }}>队友</Button>,
                        children: <>
                            <Flex gap={"small"}>
                                <ProForm.Item
                                    name={"tmp_proxy1"}
                                    label={"队友1"}
                                    labelCol={{span: 7}}
                                    style={{width: '100%'}}
                                >
                                    <Select
                                        options={proxys.map((item) => ({value: item}))}
                                    />
                                </ProForm.Item>
                                <ProForm.Item
                                    name={"tmp_proxy1_star"}
                                >
                                    <Select options={
                                        [0, 1,2,3,4,5,6]
                                            .map((item) => ({value: item}))}/>
                                </ProForm.Item>
                            </Flex>

                            <Flex gap={"small"} >
                                <ProForm.Item
                                    name={"tmp_proxy1_engine"}
                                    label={"队友1音擎"}
                                    labelCol={{span: 7}}
                                    style={{width: '100%'}}
                                >
                                    <Select
                                        options={engines.map((item) => ({value: item}))}
                                    />
                                </ProForm.Item>
                                <ProForm.Item
                                    name={"tmp_proxy1_engine_star"}
                                >
                                    <Select  options={
                                        [1,2,3,4,5]
                                            .map((item) => ({value: item}))}/>
                                </ProForm.Item>
                            </Flex>
                            <ProForm.Item
                                name={"tmp_proxy1_driver"}
                                label={"队友1驱动盘"}
                                style={{width: '100%'}}
                            >
                                <Select
                                    options={drivers.map((item) => ({value: item}))}
                                />
                            </ProForm.Item>



                            <Flex gap={"small"}>
                                <ProForm.Item
                                    name={"tmp_proxy2"}
                                    label={"队友2"}
                                    labelCol={{span: 7}}
                                    style={{width: '100%'}}
                                >
                                    <Select
                                        options={proxys.map((item) => ({value: item}))}
                                    />
                                </ProForm.Item>
                                <ProForm.Item
                                    name={"tmp_proxy2_star"}
                                >
                                    <Select options={
                                        [0, 1,2,3,4,5,6]
                                            .map((item) => ({value: item}))}/>
                                </ProForm.Item>
                            </Flex>

                            <Flex gap={"small"} >
                                <ProForm.Item
                                    name={"tmp_proxy2_engine"}
                                    label={"队友2音擎"}
                                    labelCol={{span: 7}}
                                    style={{width: '100%'}}
                                >
                                    <Select
                                        options={engines.map((item) => ({value: item}))}
                                    />
                                </ProForm.Item>
                                <ProForm.Item
                                    name={"tmp_proxy2_engine_star"}
                                >
                                    <Select  options={
                                        [1,2,3,4,5]
                                            .map((item) => ({value: item}))}/>
                                </ProForm.Item>
                            </Flex>
                            <ProForm.Item
                                name={"tmp_proxy2_driver"}
                                label={"队友2驱动盘"}
                                style={{width: '100%'}}
                            >
                                <Select
                                    options={drivers.map((item) => ({value: item}))}
                                />
                            </ProForm.Item>
                            <Button type={"primary"} onClick={testProxyDataHandle}>生成测试数据</Button>
                        </>,
                    }]}  size={'small'} bordered={false} ghost />

                </div>
                <Flex gap={'small'} style={{margin: '8px'}}>
                    <Button type="primary" htmlType="submit" style={{marginBottom: '10px'}}>开始计算</Button>
                    <Button type="primary" htmlType="button" danger onClick={()=> {
                        form.resetFields()
                        setFuzzResult(new DriverFuzzResp())
                        setBoard("")
                    }}>
                        重置
                    </Button>
                    <Button color="purple" variant="solid" onClick={() => setDriverOpen(true)}>驱动盘选择</Button>
                    <Button color="cyan" variant="solid" onClick={() => setDriverHistoryOpen(true)}>计算历史</Button>
                </Flex>
            </ProForm>
            <DriverSelectPageMemo ref={driverSelectRef} open={driverOpen} onChange={(value) => setDriverOpen(value)}/>
            <DriverHistoryPage open={driverHistoryOpen} onChange={(value) => setDriverHistoryOpen(value)} onRowDoubleClick={driverHistoryDoubleClickHandle}/>

        </Allotment.Pane>
        <div className={globalStyles['l6-split']} style={{overflow: 'auto'}}>
            <div style={{width: '98%', height: '100%'}}>
                <L6Row height={'50%'}>
                    <L6Col span={12} style={{height: '100%'}}>
                        <DriverLogPageMemo ref={driverLogRef} eventID={eventID}/>
                    </L6Col>
                    <L6Col span={12} style={{height: '100%'}}>
                        <Flex vertical style={{height: '100%'}} flex={"0 0 140px"}>
                            <span className={globalStyles['l6-label']}>面板</span>
                            <L6CodeEditor height={'calc( 100% - 32px )'} value={board} onChange={(code) => setBoard(code)}  simpleMode/>
                        </Flex>
                    </L6Col>


                </L6Row>
                <div className={globalStyles['l6-label']} style={{width: '99%'}}>驱动盘</div>
                <L6Row height={'40%'} >
                    <L6Col span={8} style={{height: '100%'}} >
                        <div className={globalStyles['l6-image']}>
                            <Image
                                src={`data:image/png;base64,${fuzzResult.disk1}`}
                                fallback="data:image/png;base64,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"
                            />
                        </div>
                    </L6Col>
                    <L6Col span={8} style={{height: '100%'}} >
                        <div className={globalStyles['l6-image']}>
                            <Image
                                src={`data:image/png;base64,${fuzzResult.disk2}`}
                                fallback="data:image/png;base64,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"
                            />
                        </div>
                    </L6Col>
                    <L6Col span={8} style={{height: '100%'}} >
                        <div className={globalStyles['l6-image']}>
                            <Image
                                src={`data:image/png;base64,${fuzzResult.disk3}`}
                                fallback="data:image/png;base64,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"
                            />
                        </div>
                    </L6Col>
                </L6Row>
                <L6Row height={'40%'} >
                    <L6Col span={8} style={{height: '100%'}} >
                        <div className={globalStyles['l6-image']}>
                            <Image
                                src={`data:image/png;base64,${fuzzResult.disk4}`}
                                fallback="data:image/png;base64,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"
                            />
                        </div>
                    </L6Col>
                    <L6Col span={8} style={{height: '100%'}} >
                        <div className={globalStyles['l6-image']}>
                            <Image
                                src={`data:image/png;base64,${fuzzResult.disk5}`}
                                fallback="data:image/png;base64,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"
                            />
                        </div>
                    </L6Col>
                    <L6Col span={8} style={{height: '100%'}} >
                        <div className={globalStyles['l6-image']}>
                            <Image
                                src={`data:image/png;base64,${fuzzResult.disk6}`}
                                fallback="data:image/png;base64,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"
                            />
                        </div>
                    </L6Col>
                </L6Row>
            </div>

        </div>
    </Allotment>
    )
}
const DriverFuzzPage: React.FC = () => {
    return (
        <PageWrapTabs
            tag={"driver_fuzz"}
            content={(id?: string)=> <Content id={id}/>}
        />
    );
}
export default DriverFuzzPage;
