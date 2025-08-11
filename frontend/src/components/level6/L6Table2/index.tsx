import 'ag-grid-enterprise';
import 'ag-grid-community/styles/ag-grid.css';

// import 'ag-grid-community/styles/ag-theme-quartz.css';
import 'ag-grid-enterprise';
import { AgGridReact } from 'ag-grid-react';

import {ExtractInfo, OperationSupport} from '@/components/level6/L6Table2/types';
import { CommonConfigOption } from '@/models/interface';

import L6DrawerForm, {L6FormFilterRules} from '@/components/level6/L6DrawerForm';
import { DynamicHandleParam } from '@/components/level6/L6DrawerForm/types';
import {compileExpression, evaluateExpression, getColumnTypeFromConfig} from '@/components/level6/L6Table2/helper';
import { MenuItemConfig } from '@/models/interface';
import { uuid } from '@/utils/uuid';
import {ClearOutlined, FilterOutlined, PlusOutlined, RedoOutlined} from '@ant-design/icons';
import {ICellRendererParams, IServerSideDatasource, IServerSideGetRowsParams} from 'ag-grid-community';
import { CellStyleFunc } from 'ag-grid-community';
import { RowDoubleClickedEvent } from 'ag-grid-community';
import { LicenseManager } from 'ag-grid-enterprise';
import {Button, Descriptions, Drawer, Flex, message, Popconfirm, Switch, Tooltip} from 'antd';
import React, { useEffect, useImperativeHandle, useMemo, useRef, useState } from 'react';
import { TriggerEvent, contextMenu } from 'react-contexify';
import './styles.less';
import {NamePath} from "rc-field-form/lib/interface";
import L6Menu from "@/components/level6/L6Menu";
import {CellDoubleClickedEvent, RowClickedEvent} from "ag-grid-community/dist/types/core/events";
import {sleep} from "@/utils/time";
import {GroupToolPanel, GroupToolPanelParams, GroupTreeNode} from "./GroupToolPanel";

const AgGridReactMemo = React.memo(AgGridReact);

LicenseManager.setLicenseKey('[v3]-[0102]_NDA3MDg4MDAwMDAwMA==e65b220ba1afe6f55b5bd31c6d22826e');




// const { confirm } = Modal;
// const showConfirm = (tag: string, okFn: ()=> void) => {
//   confirm({
//     title: `你确认要进行"${tag}"操作吗?`,
//     icon: <ExclamationCircleFilled />,
//     onOk() {
//       okFn()
//     },
//   });
// };

/*
  menu:
    const data = (menuProps as Record<string, any>).data
    const rows = (menuProps as Record<string, any>).rows
    const listData = (menuProps as Record<string, any>).listData
    const setCreateRow = (menuProps as Record<string, any>).setCreateRow
    const setCreateFormOpen = (menuProps as Record<string, any>).setCreateFormOpen

*/
export interface GroupParams {
    groupKey: string;
    treeData: GroupTreeNode[],
}
export interface L6Table2Props {
    keyName: string;
    options: CommonConfigOption[];
    dataSource: Record<string, any>[];
    operation: OperationSupport;
    manualReload?: boolean;
    title?: string;
    titleRenderPosition?: string,
    titleRender?: () => React.ReactNode[];
    formTitleRender?: (setFieldValue: (name: NamePath, value: any) => void, value?: any) => React.ReactNode;
    createFormOpenInitValue?: (setFieldValue: (name: NamePath, value: any) => void, record?: Record<string, any>) => void,
    formWidth?: number|string,

    menuOptions?: MenuItemConfig[];
    menuParams?: Record<string, any>;
    customCell?: Record<string, CellStyleFunc> | CellStyleFunc;

    onRowClick?: (event: RowClickedEvent) => Promise<void>;
    onRowDoubleClick?: (event: RowDoubleClickedEvent) => Promise<void>;
    onCellDoubleClick?:  (event: CellDoubleClickedEvent) => Promise<void>;
    onDynamicForm?: (param: DynamicHandleParam) => Promise<void>;
    icon?: (record: { [key: string]: any }) => React.ReactNode;
    showPage?: boolean;
    serverSideOption?: ServerSideOption,

    noellipsis?: boolean; // 开启后高度自适应

    suppressSideBar?: boolean;
    suppressRowSelection?: boolean; // 切换成框选删除模式
    suppressToolBar?: boolean

    rowSelection?: boolean
    filterRules?: L6FormFilterRules
    height?: string,
    enableAutoRefresh?: boolean,
    refreshInterval?: number,
    customHide?: Record<string, boolean>

    detailWidth?: number|  string
    extracts?: ExtractInfo[]

    groupParams?: GroupParams,
}
export interface ServerSideOption {
    list:  (request: IServerSideGetRowsParams) => Promise<void>
}

export type L6Table2Ref = {
    tableRef: AgGridReact | null;
    setCreateRow: (row: Record<string, any>) => void;
    setCreateFormOpen: (open: boolean) => void;

};

const L6Table2 = React.forwardRef<L6Table2Ref, L6Table2Props>((props, ref) => {
    const [showPage, setShowPage] = useState(!!props.showPage);
    const [autoRefresh, setAutoRefresh] = useState(false)

    const [editFormOpen, setEditFormOpen] = useState(false);
    const [createFormOpen, setCreateFormOpen] = useState(false);

    const [detailOpen, setDetailOpen] = useState(false);

    const [currentRow, setCurrentRow] = useState<Record<string, any>>({});
    const [createRow, setCreateRow] = useState<Record<string, any>>({});
    const tableRef = useRef<AgGridReact>(null);
    const [tableKey, setTableKey] = useState(uuid());

    const menuID = useMemo(() => uuid(), []);
    const rowCount = useMemo(() => props.dataSource.length, [props.dataSource]);

    useImperativeHandle(ref, () => ({
        // 在这里定义 ref 暴露的方法或属性
        tableRef: tableRef.current, // 将子组件内部的 tableRef 暴露给父组件
        setCreateFormOpen: setCreateFormOpen,
        setCreateRow: setCreateRow,
        // 其他方法或属性
    }));

    const columns = useMemo(() => {
        let cols = getColumnTypeFromConfig(props.options, props.customCell, props.noellipsis, props.customHide);

        if (props.icon) {
            cols = [
                {
                    field: 'icon',
                    headerName: '',
                    width: 40,
                    minWidth: 40,
                    maxWidth: 40,
                    pinned: 'left',

                    suppressMovable: true,
                    suppressHeaderMenuButton: true,
                    suppressColumnsToolPanel: true,
                    suppressFiltersToolPanel: true,

                    cellRenderer: function (params: ICellRendererParams) {
                        return props.icon ? props.icon(params.data) : <></>;
                    },
                },
                ...cols,
            ];
        }

        if ((props.rowSelection ||  props.operation.MultiDelete) && !props.suppressRowSelection) {
            cols = [
                {
                    field: 'selection',
                    headerName: '',
                    width: 40,
                    minWidth: 40,
                    maxWidth: 40,
                    pinned: 'left',

                    suppressColumnsToolPanel: true,
                    suppressHeaderMenuButton: true,
                    suppressMovable: true,
                    suppressFiltersToolPanel: true,

                    headerCheckboxSelection: true,
                    checkboxSelection: true,
                    headerCheckboxSelectionFilteredOnly: true,
                },
                ...cols,
            ];
        }

        return cols;
    }, [props.options]);
    // useEffect(() => {
    //     const keyOrder = Object.keys(currentRow); // Capture original key order
    //
    //     const a = Object.entries(currentRow)
    //         .map(([key, value]) => {
    //             const items = columns.filter(item => item.field === key);
    //             return {
    //                 key: key,
    //                 label: items.length ? items[0].headerName : key,
    //                 children: value ? value.toString() : null,
    //             };
    //         })
    //         .sort((a, b) => keyOrder.indexOf(a.key) - keyOrder.indexOf(b.key));
    //     console.dir(a)
    // }, [currentRow]);
    useEffect(() => {
        if (props.manualReload) {
            return;
        }
        props.operation.List(true);
    }, [props.operation]);
    const handleRefreshClick = async function () {
        await props.operation.List(true);
    };

    const handleResetFilterClick = async function () {
        tableRef.current?.api.setFilterModel(null)
        message.success("重置过滤成功")
    };

    const menuOptions: MenuItemConfig[] = useMemo(() => {
        const menu: MenuItemConfig[] = [];
        if (props.operation.Update) {
            menu.push({
                title: '编辑',
                key: 'edit',
                event: function (args) {
                    setEditFormOpen(true);
                },
            });
        }
        if (props.operation.Delete) {
            menu.push({
                title: '删除',
                key: 'delete',
                twiceConfirm: true,
                event: function (args) {
                    const data = (args as Record<string, any>).data;
                    // @ts-ignore
                    props.operation.Delete(data);
                },
            });
        }
        if (props.operation.MultiDelete) {
            menu.push({
                title: '删除选中',
                key: 'multi-delete',
                twiceConfirm: true,
                event: function (args) {
                    const rows = (args as Record<string, any>).rows;
                    if (rows.length === 0) {
                        return;
                    }
                    // @ts-ignore
                    props.operation.MultiDelete(rows);
                },
            });
        }

        if (props.operation.Detail) {
            menu.push({
                title: '详情',
                key: 'detail',
                event: function (args) {
                    setDetailOpen(true);
                },
            });
        }
        if (props.operation.CopyColumns) {
            menu.push({
                title: '列复制',
                key: 'copy-columns',
                event: function (args) {
                    const range = tableRef.current?.api.getCellRanges()
                    if (range && range.length) {
                        const colId = range[0].startColumn.getColId();
                        const cols: string[] = []
                        tableRef.current?.api.forEachNodeAfterFilter((rowNode, index)=> {
                            cols.push(rowNode.data[colId])
                        })

                        navigator.clipboard.writeText(cols.join("\n"))
                        message.success("复制成功")
                    }
                },
            });
        }
        if (props.extracts) {
            const subMenu: MenuItemConfig[] = [];
            for (let extract of props.extracts) {
                subMenu.push({
                    title: extract.name,
                    key: extract.name,
                    event: function (args) {
                        message.info("提取中")
                        let result: string[] = [];
                        tableRef.current?.api.forEachNodeAfterFilter((rowNode, index) => {
                            const compiledExpression = compileExpression(extract.expression); // 预编译表达式
                            let parsedResult = evaluateExpression(compiledExpression, rowNode.data);
                            if (parsedResult) {
                                if (extract.post) {
                                    parsedResult = extract.post(parsedResult)
                                }
                                result.push(parsedResult)
                            }
                        })
                        result = [...new Set(result)];
                        if (result.length) {
                            navigator.clipboard.writeText(result.join("\n"))
                            message.success(`提取到剪切板, 数量为${result.length}`)
                        } else {
                            message.warning("提取结果为空")
                        }


                        return
                    }
                })
            }
            menu.push({
                title: '信息提取',
                key: 'extract-info',
                children: subMenu
            });
        }
        if (props.operation.Exports) {
            menu.push({
                title: '导出Excel',
                key: 'exports',
                event: function (args) {
                    tableRef.current?.api.exportDataAsCsv({fileName: props.operation.Exports})
                },
            });
        }
        if (props.menuOptions) {
            if (menu.length) {
                menu.push({
                    title: '',
                    key: 'separator',
                });
            }

            menu.push(...props.menuOptions);
        }
        return menu;
    }, [props.menuOptions]);
    useEffect(() => {
        if (showPage) {
            setTableKey(uuid());
        }
    }, [showPage]);

    useEffect(() => {
        let isActive = true; // 控制循环的标志变量

        const fetch = async function() {
            if (autoRefresh) {
                while (isActive && autoRefresh){
                    await props.operation.List(false)
                    await sleep(props.refreshInterval?props.refreshInterval: 5000)
                }
            }
        }
        fetch()
        return () => {
            isActive = false;
        };
    }, [autoRefresh]);

    useEffect(() => {
        if (tableRef.current && tableRef.current.api && props.serverSideOption) {
            const  source: IServerSideDatasource = {
                getRows(params: IServerSideGetRowsParams) {
                    props.serverSideOption?.list(params)
                }
            };
            tableRef.current.api.setGridOption("serverSideDatasource", source)
        }
    }, [tableRef.current]);


    return (
        <Flex vertical style={{ width: '100%', height: props.height?props.height:'100%' }} gap={'small'}>
            {!props.suppressToolBar?<Flex
                justify={'space-between'}
                align={'center'}
                gap={'small'}
                style={{ marginRight: 20 }}
            >
                <Flex align={'center'} gap="small">
                    {props.titleRenderPosition && props.titleRenderPosition == "left"?props.titleRender ? props.titleRender() : null: null}
                    {props.operation.Create ? (
                        <Button
                            key="primary"
                            icon={<PlusOutlined />}
                            type="primary"
                            onClick={() => {
                                setCreateFormOpen(true);
                            }}
                        >
                            添加{props.title}
                        </Button>
                    ) : (
                        <span></span>
                    )}
                    {!props.titleRenderPosition || props.titleRenderPosition && props.titleRenderPosition == "right"?props.titleRender ? props.titleRender() : null: null}
                </Flex>
                <Flex justify={'flex-end'} align={'center'} gap={'small'}>
                    {props.operation.Clear ? (
                        <Popconfirm
                            title="确认清空数据吗?"
                            onConfirm={() => {
                                if (props.operation.Clear) {
                                    props.operation.Clear({})
                                }
                            }}
                        >
                            <Button
                                key="clear"
                                icon={<ClearOutlined />}
                                danger
                            >
                                清空
                            </Button>
                        </Popconfirm>

                    ) : null}
                    { props.enableAutoRefresh?
                        <Switch
                            checked={autoRefresh}
                            onChange={(checked) => setAutoRefresh(checked)}
                            checkedChildren={'自动刷新'}
                            unCheckedChildren={'停止刷新'}
                        />: null}
                    {
                        !props.serverSideOption?
                            <><Switch
                                checked={showPage}
                                onChange={(checked) => setShowPage(checked)}
                                checkedChildren={'开启分页'}
                                unCheckedChildren={'关闭分页'}
                            /><span>Total {rowCount} items</span></>:null
                    }

                    <Tooltip title="重置过滤">
                        <Button
                            shape="circle"
                            type="text"
                            size={'middle'}
                            icon={<FilterOutlined />}
                            onClick={handleResetFilterClick}
                        />
                    </Tooltip>
                    <Tooltip title="刷新">
                        <Button
                            shape="circle"
                            type="text"
                            size={'middle'}
                            icon={<RedoOutlined />}
                            onClick={handleRefreshClick}
                        />
                    </Tooltip>
                </Flex>
            </Flex>: null}

            <div
                style={{ width: '100%', height: 'calc( 100% - 32px )' }}
                className="ag-theme-balham"
            >
                <AgGridReact
                    ref={tableRef}
                    key={tableKey}
                    rowData={props.dataSource}
                    columnDefs={columns}
                    // modules={modules}
                    enableRangeSelection={true}
                    columnMenu={'new'}
                    suppressContextMenu={true}
                    rowSelection="multiple"
                    suppressRowClickSelection={true}
                    onSelectionChanged={function (event) {}}
                    enableBrowserTooltips
                    onRowClicked={props.onRowClick}
                    onRowDoubleClicked={props.onRowDoubleClick}
                    onCellDoubleClicked={props.onCellDoubleClick}
                    preventDefaultOnContextMenu={true}
                    onCellContextMenu={function (event) {
                        setCurrentRow(event.data as any);
                        if (menuOptions.length) {
                            let rows: any[] = [];
                            if (props.suppressRowSelection) {
                                const range = tableRef.current?.api.getCellRanges()?.[0];
                                if (range) {
                                    const { startRow, endRow } = range;
                                    const start = Math.min(
                                        startRow?.rowIndex as number,
                                        endRow?.rowIndex as number,
                                    );
                                    const end = Math.max(
                                        startRow?.rowIndex as number,
                                        endRow?.rowIndex as number,
                                    );

                                    for (let i = start; i <= end; i++) {
                                        const row =
                                            tableRef.current?.api.getDisplayedRowAtIndex(i)?.data;
                                        if (row) {
                                            rows.push(row);
                                        }
                                    }
                                }
                            } else {
                                rows = tableRef.current?.api.getSelectedRows() as any[];
                            }
                            let menuProps: Record<string, any> = {
                                data: event.data,
                                rows: rows,
                                setCreateRow: setCreateRow,
                                setCreateFormOpen: setCreateFormOpen,
                            };
                            if (props.menuParams) {
                                menuProps = {
                                    ...menuProps,
                                    ...props.menuParams,
                                };
                            }
                            contextMenu.show({
                                id: menuID,
                                event: event.event as TriggerEvent,
                                props: menuProps,
                            });
                        }
                    }}
                    gridOptions={{
                        icons: {
                            sortAscending:
                                '<span role="img" aria-label="caret-down" class="ag-icon" style="color: #1677ff;font-size: 11px;"><svg focusable="false" class="" data-icon="caret-up" width="1em" height="1em" fill="currentColor" aria-hidden="true" viewBox="0 0 1024 1024"><path d="M858.9 689L530.5 308.2c-9.4-10.9-27.5-10.9-37 0L165.1 689c-12.2 14.2-1.2 35 18.5 35h656.8c19.7 0 30.7-20.8 18.5-35z"></path></svg></span>',
                            sortDescending:
                                '<span role="img" aria-label="caret-down" class="ag-icon" style="color: #1677ff;font-size: 11px;"><svg focusable="false" class="" data-icon="caret-down" width="1em" height="1em" fill="currentColor" aria-hidden="true" viewBox="0 0 1024 1024"><path d="M840.4 300H183.6c-19.7 0-30.7 20.8-18.5 35l328.4 380.8c9.4 10.9 27.5 10.9 37 0L858.9 335c12.2-14.2 1.2-35-18.5-35z"></path></svg></span>',
                            filter: '<span class="ag-icon-filter" unselectable="on" role="presentation" ><svg focusable="false" class="" data-icon="search" width="1em" height="1em" fill="currentColor" aria-hidden="true" viewBox="64 64 896 896"><path d="M909.6 854.5L649.9 594.8C690.2 542.7 712 479 712 412c0-80.2-31.3-155.4-87.9-212.1-56.6-56.7-132-87.9-212.1-87.9s-155.5 31.3-212.1 87.9C143.2 256.5 112 331.8 112 412c0 80.1 31.3 155.5 87.9 212.1C256.5 680.8 331.8 712 412 712c67 0 130.6-21.8 182.7-62l259.7 259.6a8.2 8.2 0 0011.6 0l43.6-43.5a8.2 8.2 0 000-11.6zM570.4 570.4C528 612.7 471.8 636 412 636s-116-23.3-158.4-65.6C211.3 528 188 471.8 188 412s23.3-116.1 65.6-158.4C296 211.3 352.2 188 412 188s116.1 23.2 158.4 65.6S636 352.2 636 412s-23.3 116.1-65.6 158.4z"></path></svg></span>',
                        },
                    }}
                    sideBar={
                        props.suppressSideBar
                            ? null
                            : {
                                  toolPanels: props.groupParams?
                                      [
                                          {
                                              id: 'columns',
                                              labelDefault: 'Columns',
                                              labelKey: 'columns',
                                              iconKey: 'columns',
                                              toolPanel: 'agColumnsToolPanel',
                                              toolPanelParams: {
                                                  suppressPivots: true,
                                                  suppressPivotMode: true,
                                                  suppressRowGroups: true,
                                                  suppressValues: true,
                                                  suppressColumnMove: true,
                                              },
                                          },
                                          {
                                              id: "groups",
                                              labelDefault: "Groups",
                                              labelKey: "groups",
                                              toolPanel: GroupToolPanel,
                                              toolPanelParams: props.groupParams,
                                          }
                                      ]: [
                                          {
                                              id: 'columns',
                                              labelDefault: 'Columns',
                                              labelKey: 'columns',
                                              iconKey: 'columns',
                                              toolPanel: 'agColumnsToolPanel',
                                              toolPanelParams: {
                                                  suppressPivots: true,
                                                  suppressPivotMode: true,
                                                  suppressRowGroups: true,
                                                  suppressValues: true,
                                                  suppressColumnMove: true,
                                              },
                                          },
                                      ],
                                  position: 'right',
                                  defaultToolPanel: 'filters',
                                  hiddenByDefault: false,
                              }
                    }
                    pagination={showPage}
                    paginationPageSizeSelector={[20, 100, 500, 1000]}
                    paginationPageSize={20}
                    cacheBlockSize={props.serverSideOption?20: undefined}
                    rowModelType={props.serverSideOption?'serverSide': undefined}
                    autoGroupColumnDef={{
                        headerName: "分组",
                    }}
                    groupDefaultExpanded={1}
                    groupDisplayType="groupRows"
                    // groupAllowUnbalanced={true}

                />
            </div>
            {props.operation.Detail ? (
                <Drawer
                    title="详情"
                    width={props.detailWidth?props.detailWidth:500}
                    open={detailOpen}
                    onClose={function () {
                        setDetailOpen(false);
                    }}
                    styles={{
                        header: { marginTop: '16px' },
                        body: { paddingBottom: '80px' },
                        footer: { textAlign: 'right' },
                    }}
                    // getContainer={false}
                >
                    <Descriptions
                        bordered
                        column={1}
                        size={"small"}
                        contentStyle={{whiteSpace: 'pre-wrap'}}
                        labelStyle={{ backgroundColor: 'deepskyblue', whiteSpace: 'nowrap' }}
                        items={currentRow?columns.map(column => {
                            const { field, headerName } = column;
                            const value = currentRow[field as string];
                            return {
                                key: field,
                                label: headerName || field,
                                children: value ? value.toString() : null,
                            };
                        }):[]}
                    ></Descriptions>
                </Drawer>
            ) : null}
            {menuOptions.length ? <L6Menu menuID={menuID} options={menuOptions}></L6Menu> : null}
            {props.operation.Create ? (
                <L6DrawerForm
                    title={'新建' + props.title}
                    open={createFormOpen}
                    setOpen={setCreateFormOpen}
                    options={props.options}
                    record={createRow}
                    onSubmit={props.operation.Create}
                    onDynamicForm={props.onDynamicForm}
                    titleRender={props.formTitleRender}
                    openInitValue={props.createFormOpenInitValue}
                    filterRules={props.filterRules}
                    width={props.formWidth}
                />
            ) : null}

            {props.operation.Update && editFormOpen ? (
                <L6DrawerForm
                    title={'编辑' + props.title}
                    open={editFormOpen}
                    setOpen={setEditFormOpen}
                    options={props.options}
                    onSubmit={props.operation.Update}
                    record={currentRow}
                    destroyOnClose
                    onDynamicForm={props.onDynamicForm}
                    titleRender={props.formTitleRender}
                    filterRules={props.filterRules}
                    width={props.formWidth}
                />
            ) : null}
        </Flex>
    );
});

L6Table2.defaultProps = {
    key: 'uuid',
};

export default L6Table2;
