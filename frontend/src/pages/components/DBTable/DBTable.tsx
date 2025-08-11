import {DynamicHandleParam} from '@/components/level6/L6DrawerForm/types';
import {OperationSupport} from '@/components/level6/L6Table2/types';
import L6Table2, {L6Table2Props, L6Table2Ref} from '@/components/level6/L6Table2';
import {CommonConfigOption, MenuItemConfig} from '@/models/interface';
import {updateGenerateForm} from './helper';
import {CellStyleFunc, RowDoubleClickedEvent} from 'ag-grid-community';
import {message} from 'antd';
import React, {RefObject, useEffect, useImperativeHandle, useMemo, useRef, useState} from 'react';
import {Updater, useImmer} from 'use-immer';

import type {DBTableOperationSupport, DynamicValueParam, UpdateDataType} from './types';
import {ClientDBManage, GetClientConfigOptionsWithName} from '../../../../wailsjs/go/main_control/Control';
import {common_model} from "../../../../wailsjs/go/models";


/*
  menu:
    const data = (menuProps as Record<string, any>).data
    const listData = (menuProps as Record<string, any>).listData
    const setCreateRow = (menuProps as Record<string, any>).setCreateRow
    const setCreateFormOpen = (menuProps as Record<string, any>).setCreateFormOpen
*/

interface DBTableProps {
    moduleName: string;
    keyName: string;
    operation: DBTableOperationSupport;
    title?: string;
    onListData?: (updateDataSource: (data: Record<string, any>[]) => void) => Promise<boolean>;
    onRowDoubleClick?: (currentRow: Record<string, any>, update: UpdateDataType) => Promise<void>;
    onRowClick?: (currentRow: Record<string, any>, update: UpdateDataType) => Promise<void>;

    dynamicMap?: Record<string, DynamicValueParam>;
    dbCondition?: () => Record<string, any>;

    l6Props?: Partial<L6Table2Props>;

}

export type DBTableRef = {
    listData: () => Promise<boolean>;
    updateDataSource: Updater<Record<string, any>[]>
    l6TableRef: RefObject<L6Table2Ref>
};

const DBTable = React.forwardRef<DBTableRef, DBTableProps>((props, ref) => {
    const [columns, setColumns] = useState<CommonConfigOption[]>([]);
    const [dataSource, updateDataSource] = useImmer<Record<string, any>[]>([]);
    const tableRef = useRef<L6Table2Ref>(null);
    useImperativeHandle(ref, () => ({
        // 在这里定义 ref 暴露的方法或属性
        listData: listData, // 将子组件内部的 tableRef 暴露给父组件
        updateDataSource: updateDataSource,
        l6TableRef: tableRef,
        // 其他方法或属性
    }));

    const listData = props.onListData
        ? async () => {
            // @ts-ignore
            return await props.onListData(updateDataSource);
        }
        : async (): Promise<boolean> => {
            const resp = await ClientDBManage(
                new common_model.DBManageReq({
                    module: props.moduleName,
                    operation: 'list',
                    conditions: props.dbCondition ? [{ info: props.dbCondition() }] : [],
                }),
            );
            if (!resp.status) {
                message.error('刷新失败, ' + resp.err);
                return false;
            }
            const data = resp.infos.map((item) => item.info);
            updateDataSource(data);
            return true;
        };

    const listDataWithTip = async function (tips?: boolean) {
        const result = await listData();
        if (result && tips) {
            message.success('刷新成功');

        }
    };

    const updateData = props.operation.Update
        ? async (data: Record<string, any>) => {
            const condition: Record<string, any> = {};
            condition[props.keyName] = data[props.keyName];
            const resp = await ClientDBManage(
                new common_model.DBManageReq({
                    module: props.moduleName,
                    operation: 'update',
                    info: data,
                    conditions: [
                        {
                            info: props.dbCondition
                                ? {...props.dbCondition(), ...condition}
                                : condition,
                        },
                    ],
                }),
            );
            if (!resp.status) {
                message.error('更新失败, ' + resp.err);
                return;
            }
            message.success('更新成功');
            await listData();
        }
        : undefined;

    const createData = props.operation.Create
        ? async (data: Record<string, any>) => {
            const resp = await ClientDBManage(
                new common_model.DBManageReq({
                    module: props.moduleName,
                    operation: 'create',
                    info: props.dbCondition ? {...props.dbCondition(), ...data} : data,
                }),
            );
            if (!resp.status) {
                message.error('创建失败, ' + resp.err);
                return;
            }
            message.success('创建成功');
            await listData();
        }
        : undefined;

    const deleteData = props.operation.Delete
        ? async (data: Record<string, any>) => {
            const condition: Record<string, any> = {};
            condition[props.keyName] = data[props.keyName];
            const resp = await ClientDBManage(
                new common_model.DBManageReq({
                    module: props.moduleName,
                    operation: 'delete',
                    conditions: [
                        {
                            info: props.dbCondition
                                ? {...props.dbCondition(), ...condition}
                                : condition,
                        },
                    ],
                }),
            );
            if (!resp.status) {
                message.error('删除失败, ' + resp.err);
                return;
            }
            message.success('删除成功');
            await listData();
        }
        : undefined;
    const clearData = props.operation.Clear
        ? async () => {
            const condition: Record<string, any> = {};
            const resp = await ClientDBManage(
                new common_model.DBManageReq({
                    module: props.moduleName,
                    operation: 'delete',
                    conditions: [{info: props.dbCondition ? props.dbCondition() : {}}],
                }),
            );
            if (!resp.status) {
                message.error('清空失败, ' + resp.err);
                return;
            }
            message.success('清空成功');
            await listData();
        }
        : undefined;

    const multiDeleteData = props.operation.MultiDelete
        ? async (datas: Record<string, any>[]) => {
            const conditions = datas.map((item) => {
                const condition: Record<string, any> = {};
                condition[props.keyName] = item[props.keyName];
                return {
                    info: props.dbCondition
                        ? {...props.dbCondition(), ...condition}
                        : condition,
                };
            });
            const resp = await ClientDBManage(
                new common_model.DBManageReq({
                    module: props.moduleName,
                    operation: 'delete',
                    conditions: conditions,
                }),
            );
            if (!resp.status) {
                message.error('删除失败, ' + resp.err);
                return;
            }
            message.success('删除成功');
            await listData();
        }
        : undefined;

    const operation: OperationSupport = useMemo(
        () => {
            return {
                List: listDataWithTip,
                Create: createData,
                Update: updateData,
                Delete: deleteData,
                MultiDelete: multiDeleteData,
                Clear: clearData,
                Detail: props.operation.Detail,
                CopyColumns: props.operation.CopyColumns,
                Exports: props.operation.Exports? props.moduleName:undefined,
            }
        }, [props.dbCondition]);

    const handleDynamicForm = async function (param: DynamicHandleParam) {
        if (props.dynamicMap && param.key in props.dynamicMap) {
            const valueParam = props.dynamicMap[param.key];
            if (valueParam.handle) {
                const status = await valueParam.handle(
                    param.getFieldValue,
                    param.setFieldValue,
                    param.updateOptions,
                );
                if (!status) {
                    return;
                }
            } else {
                await updateGenerateForm({
                    ...valueParam,
                    ...param,
                });
            }

            if (valueParam.nextTrick) {
                await handleDynamicForm({
                    ...param,
                    key: valueParam.nextTrick,
                });
            }
        }
    };

    useEffect(() => {
        const fetchData = async () => {
            let resp: common_model.CommonConfigOptionsResp;
            resp = await GetClientConfigOptionsWithName(props.moduleName);
            if (resp.status) {
                setColumns(resp.options);
            } else {
                message.error(resp.err);
            }
        };
        fetchData();
    }, []);
    const handleDbClick = props.onRowDoubleClick
        ? async (event: RowDoubleClickedEvent) => {
            if (props.onRowDoubleClick) {
                await props.onRowDoubleClick(event.data, updateDataSource);
            }
        }
        : undefined;
    const handleClick = props.onRowClick
        ? async (event: RowDoubleClickedEvent) => {
            if (props.onRowClick) {
                await props.onRowClick(event.data, updateDataSource);
            }
        }
        : undefined;

    return (
        <L6Table2
            {...props.l6Props}
            ref={tableRef}
            title={props.title}
            keyName={props.keyName}
            options={columns}
            dataSource={dataSource}
            operation={operation}
            onRowClick={handleClick}
            onRowDoubleClick={handleDbClick}
            menuParams={{
                ...props.l6Props?.menuParams,
                ...{
                    updateDataSource: updateDataSource,
                    listData: listData,
                }
            }}
            onDynamicForm={handleDynamicForm}
        ></L6Table2>
    );
});

export default DBTable;
