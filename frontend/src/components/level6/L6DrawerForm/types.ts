import { NamePath } from 'rc-field-form/lib/interface';
import { Updater } from 'use-immer';
import { CommonConfigOption } from '@/models/interface';

export enum FormComponentType {
    Edit = 1,
    Combobox,
    TextArea,
    CheckBox,
    CheckEdit,
    CodeEditor,
    Process,
    TreeSelect,
    File,
    Password,
}

// 表单分类
export enum FormConfigType {
    BaseConfig = '基础配置',
    PayloadConfig = 'Payload上传配置',
    ReqConfig = '请求配置',
    OtherConfig = '其他配置',
    GenerateConfig = '生成Shell',
    C2ProfileConfig = 'C2Profile',
}

export type DynamicHandleParam = {
    key: string;
    getFieldValue: (name: NamePath) => any;
    setFieldValue: (name: NamePath, value: any) => void;
    updateOptions: Updater<CommonConfigOption[]>;
};
