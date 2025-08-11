import { NamePath } from 'rc-field-form/lib/interface';
import { Updater } from 'use-immer';
import {CommonConfigOption} from "@/models/interface";
import {common_model} from "../../../../wailsjs/go/models";

export interface generateParam {
    trickPrefix?: string;
    trickKey: string;

    changePrefix?: string;
    changeKey: string;

    getFieldValue: (name: NamePath) => any;
    updateOptions: Updater<CommonConfigOption[]>;
    setFieldValue: (name: NamePath, value: any) => void;
}

export interface DynamicValueParam {
    trickPrefix?: string;
    trickKey: string;

    changePrefix?: string;
    changeKey: string;
    handle?: (
        getFieldValue: (name: NamePath) => any,
        setFieldValue: (name: NamePath, value: any) => void,
        updateOptions: Updater<CommonConfigOption[]>,
    ) => Promise<boolean>;

    nextTrick?: string;
}

export type UpdateDataType = (data: Record<string, any>[]) => void;

export type DBManageFunc = (arg1: common_model.DBManageReq) => Promise<common_model.DBManageResp>;

export interface DBTableOperationSupport {
    Create?: boolean;
    Update?: boolean;
    Query?: boolean;
    Delete?: boolean;
    MultiDelete?: boolean;
    Clear?: boolean;
    Detail?: boolean;
    CopyColumns?: boolean;
    Exports?: boolean
}
