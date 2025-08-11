import {MenuItemConfig} from '@/models/interface';
import {Popconfirm} from 'antd';
import React, {useEffect, useMemo, useRef, useState} from 'react';
import { Item, ItemParams, Menu, Separator, Submenu } from 'react-contexify';
import 'react-contexify/ReactContexify.css';
import { useImmer } from 'use-immer';
import type { L6MenuComponentProps } from './components/ComponentLoader';
import ComponentLoader from './components/ComponentLoader';
import TwiceConfirmDialog from './components/TwiceConfirmDialog';
import usePopconfirm, {ClientPosition} from "@/components/level6/L6PPopconfirm";
import {uuid} from "@/utils/uuid";

type L6MenuProps = {
    options: MenuItemConfig[];
    menuID: string;
};

type MenuState = {
    state: boolean;
    props?: any;
};




const L6Menu: React.FC<L6MenuProps> = (props) => {
    const [menuStates, updateMenuStates] = useImmer<Record<string, MenuState>>({});
    // const [menus,setMenus] = useState<React.ReactNode[]|null>()
    // const [menuComponents,setMenuComponents] = useState<React.ReactNode[]|null>()
    const { showPopconfirm, popconfirm } = usePopconfirm();

    const setMenuState = function (key: string, state: MenuState) {
        updateMenuStates((draft) => {
            draft[key] = state;
            return draft;
        });
    };
    const initMenuState = function (options: MenuItemConfig[]) {
        options.map((option) => {
            if (option.children) {
                initMenuState(option.children);
            } else {
                setMenuState(option.key, { state: false });
            }
        });
    };

    const generateMenu = function (options: MenuItemConfig[], showConfirm: (newTitle: string, eventFn: () => void,pos: ClientPosition) => void) {
        return options.map(function (option) {
            if (option.children) {
                return (
                    <Submenu
                        key={option.key}
                        id={option.key}
                        label={option.title}
                        style={{ minWidth: '10px' }}
                    >
                        {generateMenu(option.children,showConfirm)}
                    </Submenu>
                );
            }
            if (option.key === 'separator') {
                return <Separator key={`${option.key}-${uuid()}`} />;
            }
            return (
                <Item
                    key={option.key}
                    id={option.key}
                    onClick={function (args: ItemParams) {
                        if (option.event) {
                            if (option.component) {
                                setMenuState(option.key, {
                                    state: true,
                                    props: args.props,
                                });
                            } else {
                                if (option.twiceConfirm) {
                                    showConfirm(option.title, function () {
                                        option.event?.(args.props, {});
                                    }, {clientX: args.triggerEvent.clientX, clientY: args.triggerEvent.clientY});
                                } else {
                                    option.event(args.props, {});
                                }
                            }
                        }
                    }}
                >
                    {option.title}
                </Item>
            );
        });
    };
    const generateMenuComponent = function (options: MenuItemConfig[]): React.ReactNode[] | null {
        return options.map(function (option) {
            if (option.children) {
                return generateMenuComponent(option.children);
            }
            return option.component && menuStates[option.key] && menuStates[option.key].state ? (
                <ComponentLoader
                    key={option.key}
                    componentFn={option.component}
                    menuProps={{
                        text: option.title,
                        open: menuStates[option.key].state,
                        callbackEvent: async function (extParams) {
                            if (option.event) {
                                option.event(menuStates[option.key].props, extParams);
                            }
                        },
                        close: function () {
                            setMenuState(option.key, { state: false });
                        },
                        rowArgs: menuStates[option.key].props,
                    }}
                />
            ) : null;
        });
    };
    // 在组件加载时调用 initMenuState
    useEffect(() => {
        initMenuState(props.options);
    }, [props.options]); // 仅在 props.options 改变时调用
    const menus = useMemo(() => generateMenu(props.options,showPopconfirm), [props.options]);
    const menuComponents = useMemo(() => generateMenuComponent(props.options), [menuStates]);

    return (
        <>
            {popconfirm}
            <Menu id={props.menuID} style={{ minWidth: '10px' }}>
                {menus}
            </Menu>
            {menuComponents}
        </>
    );
};
export default L6Menu;
export { ComponentLoader, L6MenuComponentProps, TwiceConfirmDialog };
