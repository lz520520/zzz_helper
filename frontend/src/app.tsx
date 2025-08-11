import  L6Menu from '@/components/level6/L6Menu';
import { MenuItemConfig } from '@/models/interface';
import {
    ClearOutlined,
    CustomerServiceOutlined,
    DingtalkOutlined,
    RedoOutlined,
    SyncOutlined,
} from '@ant-design/icons';
import type { Settings as LayoutSettings, RouteContextType } from '@ant-design/pro-components';
import { RouteContext, SettingDrawer } from '@ant-design/pro-components';
import type { MenuDataItem } from '@ant-design/pro-layout/es/typing';
import type { RunTimeLayoutConfig } from '@umijs/max';
import { history } from '@umijs/max';
import 'allotment/dist/style.css';
import {FloatButton, Popconfirm, message, Space, Flex} from 'antd';
import Tabs from 'antd/lib/tabs';
import { HoxRoot } from 'hox';
import React, { useEffect, useMemo, useState } from 'react';
import { contextMenu } from 'react-contexify';
import defaultSettings from '../config/defaultSettings';
import { GC, Restart } from '../wailsjs/go/main/App';
import {EventsOff, EventsOn, WindowReload} from '../wailsjs/runtime';
import styles from './utils/style.module.less';
const isDev = process.env.NODE_ENV === 'development';
const loginPath = '/user/login';
const tabTitles: Record<string, string> = {};
import globalStyles from '@/utils/style.module.less';
import {useLocation} from "@@/exports";
import {main} from "../wailsjs/go/models";

/**
 * @see  https://umijs.org/zh-CN/plugins/plugin-initial-state
 * */
export async function getInitialState(): Promise<{
    settings?: Partial<LayoutSettings>;
    currentUser?: API.CurrentUser;
    loading?: boolean;
    fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
    return {
        settings: defaultSettings as Partial<LayoutSettings>,
    };
}

interface TabItem {
    id: string;
    pathname: string;
}



// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {

    const [activeTab, setActiveTab] = useState<string>('');
    const [tabItems, setTabItems] = useState<TabItem[]>([]);

    // const [appStatus, setAppStatus] = useState<main.AppStatus>(new main.AppStatus())

    const getCurrTab = (newActiveTab: string) => tabItems.find((item) => item.id === newActiveTab);
    // 切换 Tab
    const switchTab = (newActiveTab: string) => {
        const currTab = getCurrTab(newActiveTab);
        if (currTab) {
            history.push(currTab.pathname);
            setActiveTab(newActiveTab);
        }
    };
    // 激活 Tab
    const activateTab = () => {
        const { location } = history;
        const currTab: any = tabItems.find((item) => item.pathname === location.pathname);
        if (currTab) {
            setActiveTab(currTab.id);
        }
    };

    // 移除 Tab
    const removeTab = (tabKey: string) => {
        let newActiveTab = activeTab;
        let lastIndex = -1;
        tabItems.forEach((item, i) => {
            if (item.id === tabKey) {
                lastIndex = i - 1;
            }
        });
        const newPanes = tabItems.filter((item) => item.id !== tabKey);
        if (newPanes.length && newActiveTab === tabKey) {
            if (lastIndex >= 0) {
                newActiveTab = newPanes[lastIndex].id;
            } else {
                newActiveTab = newPanes[0].id;
            }
        }
        setTabItems(newPanes);
        switchTab(newActiveTab);
    };
    // 任何 Tab 变动，激活正确的 Tab，并更新缓存
    useEffect(() => {
        activateTab();
        // localStorage.setItem('l6TabPages', JSON.stringify(tabItems));
    }, [tabItems]);

    useEffect(() => {
        if (history.location.pathname === "/" || history.location.pathname === '') {
            // @ts-ignore
            // history.push('/' + initialState.settings.app)
        }
    }, []);

    const menuOptions: MenuItemConfig[] = useMemo(
        () => [
            {
                title: '关闭其他标签',
                key: 'close-other-all',
                event: function (rowArgs) {
                    const newPanes = tabItems.filter((item) => item.id === rowArgs);
                    setTabItems(newPanes);
                },
            },
            {
                title: '关闭所有标签',
                key: 'close-all',
                event: function (rowArgs) {
                    setTabItems([]);
                },
            },
        ],
        [tabItems],
    );

    return {
        onPageChange: () => {
            const { location, action } = history;
            // 如果没有登录，重定向到 login
            // if (!initialState?.currentUser && location.pathname !== loginPath) {
            //   // history.push(loginPath);
            // }
            const pathname = location.pathname;
            const currtabItem = {
                id: location.key,
                pathname,
            };
            // @ts-ignore
            if (pathname !== '/' && pathname !== '/' + initialState.settings.app) {
                // 构建开启的 Tab 列表，并更新国际化的 Tab 标题
                setTabItems((prev: TabItem[]) => {
                    const next = prev.find((item) => item.pathname === pathname)
                        ? prev
                        : [...prev, currtabItem];
                    return next.map((item) => ({ ...item, title: tabTitles[item.pathname] }));
                });
            }
            activateTab();
        },
        menuFooterRender:(props) => {
            const [appStatus, setAppStatus] = useState<main.AppStatus>(new main.AppStatus())
            useEffect(() => {
                EventsOn("app_status", (e: main.AppStatus) => {
                    setAppStatus(e)
                })
                return () => {
                    EventsOff('app_status')
                }
            }, []);

            return (<Flex gap={"small"} vertical>
                <span>CPU: {appStatus.cpu_percent}</span>
                <span>Mem: {appStatus.memory_usage} </span>
            </Flex>);
        },
        menuHeaderRender: undefined,
        menuRender: (props, defaultDom) => {
            return <div className={globalStyles["l6-side-menu"]}>{defaultDom}</div>
        },
        // 增加一个 loading 的状态
        childrenRender: (children) => {
            const [hiddenTab, setHiddenTab] = useState(true);

            const handleSwitchTab = function () {
                setHiddenTab(!hiddenTab);
            };
            const refreshPage = function () {
                WindowReload();
            };
            const gcFunc = async function () {
                await GC();
                message.success('gc success');
            };
            const manualFunc = async function () {
                const linkHref =
                    'https://www.yuque.com/';
                window.open(linkHref, '_blank');
            };
            const restartFunc = async function () {
                const resp = await Restart();
                if (!resp.status) {
                    message.error(resp.err);
                }
            };

            // if (initialState?.loading) return <PageLoading />;
            return (
                <>
                    <RouteContext.Consumer>
                        {(ctx: RouteContextType) => {
                            // 从上下文的 routes 中构建 Map，引用各页面的 children
                            const tabContents: Record<string, any> = {};
                            const getTabContents = (
                                parentPath: string,
                                arr: MenuDataItem[] = [],
                            ) => {
                                arr.forEach((ele) => {
                                    const path =
                                        '/' + `${parentPath}/${ele.path}`.replace(/^\/+/, '');

                                    tabContents[path] = ele.element;
                                    if (ele.children) {
                                        getTabContents(path, ele.children);
                                    }
                                });
                            };
                            // 从上下文构建 Map，缓存国际化的 Tab 标题
                            const getTabTitles = (arr: MenuDataItem[] = []) => {
                                arr.forEach((ele) => {
                                    if (ele.name) {
                                        tabTitles[ele.path as string] = ele.name;
                                    }
                                    if (ele.children) {
                                        getTabTitles(ele.children);
                                    }
                                });
                            };

                            if ((ctx as Record<any, any>).route) {
                                getTabContents('', (ctx as Record<any, any>).route.routes);
                                getTabTitles(ctx.menuData);
                            }


                            return (
                                <HoxRoot>
                                    <Tabs
                                        className={styles['l6-tabs']}
                                        type="editable-card"
                                        hideAdd
                                        style={{ height: '95vh' }}
                                        onChange={switchTab}
                                        activeKey={activeTab}
                                        tabBarStyle={hiddenTab ? { display: 'none' } : {}}
                                        onEdit={function (e, action) {
                                            removeTab(e as string);
                                        }}
                                        items={tabItems.map((item) => {
                                            return {
                                                key: item.id,
                                                label: tabTitles[item.pathname],
                                                children: tabContents[item.pathname],
                                            };
                                        })}
                                        renderTabBar={(props, DefaultTabBar) => (
                                            <div
                                                onContextMenu={(event) => {
                                                    contextMenu.show({
                                                        id: 'menu-tabs',
                                                        event: event,
                                                        props: activeTab,
                                                    });
                                                }}
                                            >
                                                <DefaultTabBar {...props} />
                                                <L6Menu
                                                    options={menuOptions}
                                                    menuID={'menu-tabs'}
                                                />
                                            </div>
                                        )}
                                    ></Tabs>
                                    <FloatButton.Group
                                        trigger={'click'}
                                        type={'primary'}
                                        style={{ right: '16px', bottom: '16px' }}
                                        icon={<CustomerServiceOutlined />}
                                    >
                                        <FloatButton
                                            tooltip="刷新"
                                            icon={<SyncOutlined />}
                                            onClick={refreshPage}
                                        />
                                        <FloatButton
                                            tooltip="GC"
                                            icon={<ClearOutlined />}
                                            onClick={gcFunc}
                                        />
                                        <FloatButton
                                            tooltip="使用手册"
                                            icon={<DingtalkOutlined />}
                                            onClick={manualFunc}
                                        />
                                        <Popconfirm
                                            title={'你确定要重启吗？'}
                                            onConfirm={restartFunc}
                                        >
                                            <FloatButton tooltip={'重启'} icon={<RedoOutlined />} />
                                        </Popconfirm>
                                        <FloatButton tooltip="tab切换" onClick={handleSwitchTab} />
                                    </FloatButton.Group>
                                </HoxRoot>
                            );
                        }}
                    </RouteContext.Consumer>
                    {isDev && (
                        <SettingDrawer
                            disableUrlParams
                            enableDarkTheme
                            settings={initialState?.settings}
                            onSettingChange={(settings) => {
                                setInitialState((preInitialState) => ({
                                    ...preInitialState,
                                    settings,
                                }));
                            }}
                        />
                    )}
                </>
            );
        },
        ...initialState?.settings,
    };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
// export const request = {
//   ...errorConfig,
// };
