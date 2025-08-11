import { L6MenuComponentProps } from '@/components/level6/L6Menu';
import React from 'react';

export type ClickEvent = (
    menuProps: Record<string, any> | string,
    extParams: Record<string, any>,
) => void;
// export type MultiClickEvent = (tableArgs: Record<string, any>[], extParams: Record<string, any>) => void;

export interface MenuItemConfig {
    title: string;
    key: string;

    twiceConfirm?: boolean;
    twiceTipKey?: string;

    component?: React.FC<L6MenuComponentProps>;

    disabled?: boolean;
    event?: ClickEvent;
    // multiEvent?: MultiClickEvent,

    children?: MenuItemConfig[];
}

export interface CommonConfigOption {
    title: string;
    key: string;
    web_key: string;
    default_value: any;
    default_value_dynamic?: boolean;
    tips?: string;
    default_options?: string[];
    auto_setting?: boolean;
    form_required?: boolean;
    form_hidden?: boolean,
    form_component_type?: number;
    form_config_type?: string;
    show_in_table: boolean;
    edit?: boolean;
    sort?: string;
    col_width?: number;
    col_fixed?: string;
    col_hidden?: boolean;
    custom_filter_dropdown?: boolean;
    row_group?: boolean,
}
