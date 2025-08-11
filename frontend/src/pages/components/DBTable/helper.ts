import { generateParam } from './types';
import { message } from 'antd';
import {DynamicFormGenerate} from "../../../../wailsjs/go/main_control/Control";
import {common_model} from "../../../../wailsjs/go/models";

const comboboxExtractValue = function (str: string) {
    const match = str.match(/^(.*?)(-{{(.*?)}})?$/);
    if (match && match[1]) {
        return match[1];
    }
    return str;
};
export const updateGenerateForm = async function (param: generateParam) {
    let trickPrefix = '';
    if (param.trickPrefix) {
        const m = param.trickPrefix.match(/\{\{(.+?)\}\}/);
        if (m) {
            trickPrefix = param.trickPrefix.replace(m[0], param.getFieldValue(m[1]));
        } else {
            trickPrefix = param.trickPrefix;
        }
    }
    const resp = await DynamicFormGenerate(
        new common_model.DynamicFormReq({
            change_key: (param.changePrefix ? param.changePrefix : '') + param.changeKey,
            trick_key: trickPrefix + param.getFieldValue(param.trickKey),
        }),
    );
    if (!resp.status) {
        message.error(resp.err);
        return;
    }
    param.updateOptions((draft) => {
        draft.map((item) => {
            if (item.key === param.changeKey) {
                item.default_value = resp.values;
            }
        });
    });
    if (resp.values.length > 0) {
        // message.info(param.form.value[param.changeKey])
        const oldValue = param.getFieldValue(param.changeKey);
        const filterValues = resp.values.filter((item) => {
            return comboboxExtractValue(item) === oldValue;
        });
        if (filterValues.length == 0) {
            param.setFieldValue(param.changeKey, comboboxExtractValue(resp.values[0]));
        }
    } else {
        param.setFieldValue(param.changeKey, '');
    }
};
