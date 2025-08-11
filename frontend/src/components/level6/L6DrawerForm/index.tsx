import L6CodeEditor from '@/components/level6/L6CodeEditor';
import {DynamicHandleParam, FormComponentType} from '@/components/level6/L6DrawerForm/types';
import L6SearchInput from '@/components/level6/L6SearchInput';
import globalStyles from '@/utils/style.module.less';
import {
    DrawerForm,
    ProFormCheckbox,
    ProFormInstance,
    ProFormSelect,
    ProFormText,
    ProFormTextArea,
} from '@ant-design/pro-components';
import {ProForm} from '@ant-design/pro-form';
import {Button, Flex, Form, Input, message, Tabs} from 'antd';
import {cloneDeep} from 'lodash';
import {NamePath} from 'rc-field-form/lib/interface';
import React, {useEffect, useImperativeHandle, useMemo, useRef, useState} from 'react';
import {useImmer} from 'use-immer';
import {CommonConfigOption} from '@/models/interface';
import L6CheckInput from "@/components/level6/L6CheckInput";
import {fileOpen} from "browser-fs-access";
import {WrapFileWithHandle} from "@/utils/wrap_file_with_handle";


const comboboxExtractValue = function (str: string) {
    const match = str.match(/^(.*?)(-{{(.*?)}})?$/)
    if (match && match[1]) {
        return match[1]
    }
    return str
}
const comboboxExtractShow = function (str: string) {
    const match = str.match(/^(.*?)(-{{(.*?)}})?$/)
    if (match) {
        if (match[3]) {
            return match[3]
        } else if (match[1]) {
            return match[1]
        }
    }
    return str
}
function extractKeys(text: string): string[] {
    const regex = /{{(\w+)}}/g;
    let match: RegExpExecArray | null;
    const matches: string[] = [];

    // 使用循环调用regex.exec来获取所有匹配项
    while ((match = regex.exec(text))) {
        // 第一个捕获组是我们需要的值，索引为1
        matches.push(match[1]);
    }

    return matches;
}
const InputFile: React.FC<{value?: any, onChange?: (value: any) => void}> = ({value, onChange}) => {
    return <Flex gap={"small"}>
        <Input readOnly value={value? value.toString(): ""}/>
        <Button type={"primary"} onClick={async ()=> {
            if (!onChange) {
                return
            }
            const handle = await fileOpen()
            if (handle) {
                onChange(new WrapFileWithHandle(handle))
            } else {
                onChange(undefined)
            }
        }}>选择文件</Button>
        <Button onClick={() => {
                if (onChange) {
                    onChange(undefined)
                }
            }
        }>重置</Button>
    </Flex>
}

function FormField(
    item: CommonConfigOption,
    handle: (key: string) => Promise<void>,
    record?: Record<string, any>,
    setFieldValue?: (name: NamePath, value: any) => void,
    filterCheckKey?: string, // 检查的key
    setFormFilterKey?: (key: string) => void, // 设置check key对应的value
) {
    if (item.auto_setting) {
        return null;
    }
    let form: React.ReactNode = null;
    switch (item.form_component_type) {
        case FormComponentType.Edit:
            form = (
                <ProForm.Item
                    key={item.key}
                    name={item.key}
                    label={item.title}
                    tooltip={item.tips}
                    initialValue={item.default_value}
                    rules={[{ required: item.form_required }]}
                >
                    <L6SearchInput key={item.key} baseOptions={item.default_options} />
                </ProForm.Item>
            );
            break;
        case FormComponentType.Combobox:
            form = (
                <ProFormSelect
                    showSearch
                    hidden={item.form_hidden}
                    key={item.key}
                    name={item.key}
                    label={item.title}
                    tooltip={item.tips}
                    initialValue={item.default_value && Array.isArray(item.default_value)? (item.default_value as string[])[0] : ''}
                    placeholder={'请输入' + item.title}
                    options={
                        item.default_value
                            ? (item.default_value as string[]).map((item) => ({
                                  value: comboboxExtractValue(item),
                                  label: comboboxExtractShow(item),
                              }))
                            : []
                    }
                    onChange={async (value) => {
                        if (setFormFilterKey && filterCheckKey && filterCheckKey === item.key) {
                            setFormFilterKey(value as string);
                        }
                        await handle(item.key);
                    }}
                    rules={[{ required: item.form_required }]}
                />
            );
            break;
        case FormComponentType.TextArea:
            form = (
                <ProFormTextArea
                    key={item.key}
                    name={item.key}
                    label={item.title}
                    initialValue={item.default_value}
                    tooltip={item.tips}
                    placeholder={'请输入' + item.title}
                    rules={[{ required: item.form_required }]}
                    fieldProps={{
                        style: { textWrap: 'nowrap' },
                        rows: 6,
                    }}
                />
            );
            break;
        case FormComponentType.File:
            form = (
                <ProForm.Item
                    key={item.key}
                    name={item.key}
                    label={item.title}
                    tooltip={item.tips}
                    initialValue={item.default_value}
                    rules={[{ required: item.form_required }]}
                >
                    <InputFile />
                </ProForm.Item>
            );
            break;
        case FormComponentType.CheckBox:
            form = (
                <ProFormCheckbox
                    key={item.key}
                    name={item.key}
                    label={item.title}
                    initialValue={item.default_value}
                    tooltip={item.tips}
                    placeholder={'请输入' + item.title}
                    rules={[{ required: item.form_required }]}
                />
            );
            break;
        case FormComponentType.CheckEdit:
            form = (
                <ProForm.Item
                    key={item.key}
                    name={item.key}
                    label={item.title}
                    tooltip={item.tips}
                    initialValue={item.default_value}
                    rules={[{ required: item.form_required }]}
                >
                    <L6CheckInput key={item.key} />
                </ProForm.Item>
            );
            break;
        case FormComponentType.CodeEditor:
            form = (
                <ProForm.Item key={item.key} name={item.key} initialValue={item.default_value &&
                    item.default_value.startsWith("lang:")? "": item.default_value}>
                    <L6CodeEditor key={item.key} language={item.default_value &&
                    item.default_value.startsWith("lang:")? item.default_value.split("lang:")[1]: "yaml"} height="70vh" />
                </ProForm.Item>
            );
            break;
        case FormComponentType.Password:
            form = (
                <ProFormText.Password
                    key={item.key}
                    name={item.key}
                    label={item.title}
                    initialValue={item.default_value}
                    tooltip={item.tips}
                    placeholder={'请输入' + item.title}
                    rules={[{ required: item.form_required }]}
                />
            );
            break;
    }
    return form;
}
interface L6FormFilterRule {
    key: string,
    excludes?: string[];
    includes?: string[];
}
export interface L6FormFilterRules {
    checkKey: string;
    checkValueFormat: string;
    rules: L6FormFilterRule[];
}
interface L6DrawerFormProps {
    title: string;
    open: boolean;
    setOpen: (open: boolean) => void;
    options: CommonConfigOption[];
    onSubmit: (data: Record<string, any>) => Promise<void>;
    record?: Record<string, any>;
    destroyOnClose?: boolean;
    onDynamicForm?: (param: DynamicHandleParam) => Promise<void>;
    titleRender?: (setFieldValue: (name: NamePath, value: any) => void, value?: any) => React.ReactNode
    openInitValue?:(setFieldValue: (name: NamePath, value: any) => void, record?: Record<string, any>) => void,
    filterRules?: L6FormFilterRules,
    width?: number | string
}

export type L6DrawerFormRef = {
    formRef: ProFormInstance | undefined;
};
const L6DrawerForm = React.forwardRef<L6DrawerFormRef, L6DrawerFormProps>((props, ref) => {
    const [form] = Form.useForm();
    const formRef = useRef<ProFormInstance>();
    const [options, updateOptions] = useImmer<CommonConfigOption[]>([]);
    const [formFilterKey, setFormFilterKey] = useState('')
    const groups = useMemo(
        () =>
            Array.from(
                new Set(
                    props.options
                        .filter((option) =>  option.form_config_type) // 只考虑form_required为true的选项
                        .map((option) => option.form_config_type as string), // 提取form_config_type
                ),
            ),
        [props.options],
    );
    useImperativeHandle(ref, () => ({
        formRef: formRef.current,
    }));

    useEffect(() => {
        if (props.filterRules) {
            setFormFilterKey(props.filterRules?.checkKey)
        }
    }, [props.filterRules])

    useEffect(() => {
        updateOptions(props.options.map((item) => cloneDeep(item)));
    }, [props.options]);

    useEffect(() => {
        const fetchData = async function () {

            for (let key in props.record) {
                formRef.current?.setFieldValue(key, props.record[key]);
            }
        };
        if (props.record && formRef.current) {
            fetchData();
        }
    }, [props.record,formRef.current]);

    const handleDynamicForm = async function (key: string) {
        if (props.onDynamicForm && formRef.current) {
            await props.onDynamicForm({
                key: key,
                setFieldValue: formRef.current.setFieldValue,
                getFieldValue: formRef.current.getFieldValue,
                updateOptions: updateOptions,
            });

        }
    };
    useEffect(function () {
            const fetch = async function() {
                if (props.open) {
                    await handleDynamicForm('init')
                    setTitle(props.title)

                    if (props.openInitValue && formRef.current) {
                        props.openInitValue(formRef.current.setFieldValue, formRef.current?.getFieldsValue())
                    }

                }
            }
            fetch()

        },
        [props.open,formRef.current],
    );
    const useOptions = useMemo(
        () => {
            let finalFilterKey = "";
            if (props.filterRules && formRef.current) {
                finalFilterKey = props.filterRules.checkValueFormat;
                const keys = extractKeys(props.filterRules.checkValueFormat)
                for (let key of keys) {
                    const value = formRef.current?.getFieldValue(key)
                    if (value) {
                        finalFilterKey = finalFilterKey.replaceAll(`{{${key}}}`, value)
                    } else {
                        // return options
                    }
                }

            }
            return options.filter((config) => {
                if (!props.filterRules) {
                    return true
                }
                let status = false;
                if (finalFilterKey === ''  || config.key === props.filterRules?.checkKey) {
                    status = true;
                } else if (props.filterRules) {
                    for (let i = 0; i < props.filterRules.rules.length; i++) {
                        const rule = props.filterRules.rules[i]
                        if (finalFilterKey.match(rule.key)) {
                            if ( rule.excludes && rule.excludes.length > 0) {
                                // 当前过滤规则里是否包含当前循环中的config.key
                                status = !rule.excludes.includes(config.key);
                            } else {
                                status = true;
                            }
                            // 如果被excludes过滤了，status为false，则不需要includes检测
                            if (status) {
                                if (rule.includes && rule.includes.length > 0) {
                                    status = rule.includes.includes(config.key);
                                } else {
                                    status = true;
                                }
                            }
                            break
                        }
                    }

                } else {
                    status = true;
                }
                return status;
            })
        },
        [options, formFilterKey,formRef.current]
    );
    const [title, setTitle] = useState("")

    return (
        <DrawerForm
            formRef={formRef}
            title={
                <Flex justify={"space-between"} gap="small">
                    <div>{title}</div>
                    <Flex gap={"small"}>{props.titleRender ?
                        props.titleRender(formRef.current?.setFieldValue as (name: NamePath, value: any) => void, formRef.current?.getFieldsValue())
                        :null}</Flex>
                </Flex>
            }
            open={props.open}
            onOpenChange={props.setOpen}
            submitter={{
                render: (props, defaultDoms) => {
                    return [
                        ...defaultDoms,
                        <Button
                            key="extra-reset"
                            onClick={() => {
                                form.resetFields()
                                message.success("重置成功")
                            }}
                        >
                            重置
                        </Button>,
                    ];
                },
            }}
            form={form}
            layout="horizontal"
            labelAlign="left"
            labelCol={{ span: 6 }}
            autoFocusFirstInput
            drawerProps={{
                width: props.width?props.width: 520,
                destroyOnClose: props.destroyOnClose,
            }}
            // submitTimeout={2000}
            onFinish={async (values) => {
                const data = props.record ? { ...props.record, ...values } : values;
                await props.onSubmit(data);
                // 不返回不会关闭弹框
                return true;
            }}
        >
            <Tabs
                hideAdd
                className={globalStyles['l6-tabs']}
                items={groups.map((group) => ({
                    key: group,
                    label: group,
                    forceRender: true,
                    children: useOptions
                        .filter((item) => item.form_config_type === group)
                        .map((item) =>
                            FormField(
                                item,
                                handleDynamicForm,
                                props.record,
                                formRef.current?.setFieldValue,
                                props.filterRules?.checkKey,
                                setFormFilterKey,
                            ),
                        ),
                }))}
            />
        </DrawerForm>
    );
});

export default L6DrawerForm;
