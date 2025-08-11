import { L6MenuComponentProps } from '@/components/level6/L6Menu';
import {ProFormCheckbox, ProFormSelect, ProFormText} from '@ant-design/pro-components';
import { ProForm } from '@ant-design/pro-form';
import {Button, Flex, Form, Input, Modal} from 'antd';
import React, {useState} from 'react';
import {fileOpen} from "browser-fs-access";

import {WrapFileWithHandle} from "@/utils/wrap_file_with_handle";
import L6SearchInput from "@/components/level6/L6SearchInput";
import L6CodeEditor from "@/components/level6/L6CodeEditor";

export type L6MenuFormProps = {
    componentProps: L6MenuComponentProps;
    formOptions: L6MenuFormOption[];
    labelCol?: number;
    maskClosable?: boolean
    width?: string|number,
    resultRender?: () => React.ReactNode,

};

export interface L6MenuFormOption {
    key: string;
    title: string;
    formType: string;
    defaultValue?: any;
    options?: string[],
    tips?: string;
    required?: boolean;
    placeholder?: string;
    extParams?: Record<string, any>;
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
    </Flex>
}


const CommonFormDialog: React.FC<L6MenuFormProps> = (props) => {
    const [form] = Form.useForm();
    return (
        <Modal
            title={props.componentProps.text}
            open={props.componentProps.open}
            onCancel={props.componentProps.close}
            onOk={function () {
            }}
            footer={null}
            destroyOnClose
            maskClosable={props.maskClosable}
            width={props.width}
        >
            <div>
                <ProForm
                    layout="horizontal"
                    labelAlign="left"
                    form={form}
                    labelCol={{ span: props.labelCol }}
                    autoFocusFirstInput
                    onFinish={async (values) => {
                        await props.componentProps.callbackEvent(values);
                        if (props.componentProps.submit) {
                            props.componentProps.submit()
                        } else {
                            props.componentProps.close();
                        }
                    }}
                    submitter={{
                        render: (a, defaultDoms) => {
                            return (
                                <Flex gap={'small'} align={'center'} justify={'flex-end'}>
                                    {[
                                        ...defaultDoms,
                                        <Button
                                            key="menu-form-extra-close"
                                            onClick={props.componentProps.close}
                                        >
                                            取消
                                        </Button>,
                                    ]}
                                </Flex>
                            );
                        },
                    }}
                >
                    {props.formOptions.map((item) => {
                        let node: React.ReactNode = null;

                        let options: string[] = []
                        let initialValue = ""
                        if (item.options) {
                            options = item.options
                        } else {
                            options = item.defaultValue
                        }
                        if (item.defaultValue) {
                            if (typeof item.defaultValue === "object") {
                                initialValue = item.defaultValue[0]
                            } else {
                                initialValue = item.defaultValue
                            }
                        }

                        switch (item.formType) {
                            case 'edit':
                                node = (
                                    <ProFormText
                                        key={item.key}
                                        name={item.key}
                                        initialValue={item.defaultValue}
                                        label={item.title}
                                        tooltip={item.tips}
                                        placeholder={item.placeholder}
                                        rules={[{ required: item.required }]}
                                    />
                                );
                                break;
                            case "file":
                                node = (
                                    <ProForm.Item
                                        key={item.key}
                                        name={item.key}
                                        label={item.title}
                                        tooltip={item.tips}
                                        rules={[{ required: item.required }]}
                                    >
                                        <InputFile />
                                    </ProForm.Item>
                                )
                                break;
                            case "search":
                                node = (
                                    <ProForm.Item
                                        key={item.key}
                                        name={item.key}
                                        initialValue={initialValue}
                                        label={item.title}
                                        tooltip={item.tips}
                                        rules={[{ required: item.required }]}
                                    >
                                        <L6SearchInput key={item.key} baseOptions={options} placeholder={item.placeholder} />
                                    </ProForm.Item>
                                )
                                break;
                            case 'checkbox':
                                node = (
                                    <ProFormCheckbox
                                        key={item.key}
                                        name={item.key}
                                        label={item.title}
                                        tooltip={item.tips}
                                        initialValue={initialValue}
                                        placeholder={item.placeholder}
                                        rules={[{ required: item.required }]}
                                    />
                                );
                                break;
                            case 'combobox':
                                node = (
                                    <ProFormSelect
                                        key={item.key}
                                        name={item.key}
                                        label={item.title}
                                        tooltip={item.tips}
                                        initialValue={initialValue}
                                        placeholder={item.placeholder}
                                        options={options.map((item) => ({
                                            value: item,
                                            label: item,
                                        }))
                                        }
                                        rules={[{ required: item.required }]}
                                    />
                                );
                                break;
                            case 'code':
                                node = (
                                    <ProForm.Item key={item.key} name={item.key} initialValue={initialValue}>
                                        <L6CodeEditor key={item.key} language={item.extParams?item.extParams["lang"]: undefined}  height="60vh" />
                                    </ProForm.Item>
                                )
                                break;
                        }
                        return node;
                    })}
                </ProForm>
                {props.resultRender?props.resultRender(): undefined}
            </div>


        </Modal>
    );
};
CommonFormDialog.defaultProps = {
    labelCol: 5,
};
export default CommonFormDialog;
