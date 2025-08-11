import  L6Menu from '@/components/level6/L6Menu';
import { MenuItemConfig } from '@/models/interface';
import { uuid } from '@/utils/uuid';
import {Flex, TreeDataNode} from 'antd';
import { Input, Tree } from 'antd';
import React, { useEffect, useImperativeHandle, useMemo, useRef, useState } from 'react';
import { contextMenu } from 'react-contexify';
import './style.module.less';
import styles from './style.module.less';

const { DirectoryTree } = Tree;
const { Search } = Input;

type L6TreeProps = {
    defaultData: TreeDataNode[];
    defaultExpandAll?: boolean;
    showSearch?: boolean;
    onClick?: (e: React.MouseEvent<HTMLSpanElement>, node: TreeDataNode) => void;
    onDoubleClick?: (e: React.MouseEvent<HTMLSpanElement>, node: TreeDataNode) => void;
    menuOptions?: MenuItemConfig[];

    selectKeys?: React.Key[];
    setSelectKeys?: (keys: React.Key[]) => void;

    expandedKeys?: React.Key[];
    setExpandedKeys?: (keys: React.Key[]) => void;
    showIcon?: boolean;
};

export type L6TreeRef = {
    treeRef: any;
};

const L6Tree = React.forwardRef<L6TreeRef, L6TreeProps>((props, ref) => {
    const menuID = useMemo(() => uuid(), []);
    const treeRef = useRef<any>();
    const [internalExpandedKeys, setInternalExpandedKeys] = useState<React.Key[]>([]);
    const [internalSelectKeys, setInternalSelectKeys] = useState<React.Key[]>([]);

    const [searchValue, setSearchValue] = useState('');
    const [autoExpandParent, setAutoExpandParent] = useState(true);
    const [dataAllKeys, setDataAllKeys] = useState<{ key: React.Key; title: string }[]>([]);

    useImperativeHandle(ref, () => ({
        treeRef: treeRef.current,
    }));
    // 生成展开树
    const generateKeys = (data: TreeDataNode[]): { key: React.Key; title: string }[] => {
        const keys: { key: React.Key; title: string }[] = [];
        for (let i = 0; i < data.length; i++) {
            const node = data[i];
            const { key } = node;
            keys.push({ key, title: node.title as string });
            if (node.children) {
                keys.push(...generateKeys(node.children));
            }
        }
        return keys;
    };

    const getParentKey = (key: React.Key, tree: TreeDataNode[]): React.Key => {
        let parentKey: React.Key;
        for (let i = 0; i < tree.length; i++) {
            const node = tree[i];
            if (node.children) {
                if (node.children.some((item) => item.key === key)) {
                    parentKey = node.key;
                } else if (getParentKey(key, node.children)) {
                    parentKey = getParentKey(key, node.children);
                }
            }
        }
        return parentKey!;
    };

    const selectKeys = useMemo(() => {
        if (props.selectKeys != undefined) {
            return props.selectKeys;
        }
        return internalSelectKeys;
    }, [props.selectKeys, internalSelectKeys]);

    const setExpandedKeys = props.setExpandedKeys ? props.setExpandedKeys : setInternalExpandedKeys;

    const expandedKeys = useMemo(() => {
        if (props.expandedKeys != undefined) {
            return props.expandedKeys;
        }
        return internalExpandedKeys;
    }, [props.expandedKeys, internalExpandedKeys]);

    const setSelectKeys = props.setSelectKeys ? props.setSelectKeys : setInternalSelectKeys;

    const onExpand = (newExpandedKeys: React.Key[]) => {
        setExpandedKeys(newExpandedKeys);
        setAutoExpandParent(false);
    };
    const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = e.target;
        const newExpandedKeys = dataAllKeys
            .map((item) => {
                if (item.title.toLowerCase().indexOf(value.toLowerCase()) > -1) {
                    return getParentKey(item.key, props.defaultData);
                }
                return null;
            })
            .filter((item, i, self): item is React.Key => !!(item && self.indexOf(item) === i));


        setExpandedKeys(newExpandedKeys);
        setSearchValue(value);
        setAutoExpandParent(true);
    };
    useEffect(() => {
        const keys = generateKeys(props.defaultData);
        setDataAllKeys(keys);
        // onChange({ target: { value: '' } } as React.ChangeEvent<HTMLInputElement>);
        if (props.defaultExpandAll) {
            setExpandedKeys(keys.map(item => item.key))
        }
    }, [props.defaultData]);

    useEffect(() => {
        if (selectKeys.length > 0) {
            treeRef.current.scrollTo(selectKeys[0]);
        }
    }, [expandedKeys]);
    const treeData = useMemo(() => {
        const loop = (data: TreeDataNode[]): TreeDataNode[] =>
            data.map((item) => {
                const strTitle = item.title as string;
                const index = strTitle.toLowerCase().indexOf(searchValue.toLowerCase());
                const beforeStr = strTitle.substring(0, index);
                const searchStr = strTitle.slice(index, index + searchValue.length);
                const afterStr = strTitle.slice(index + searchValue.length);
                const title =
                    index > -1 ? (
                        <span>
                            {beforeStr}
                            <span className={styles['site-tree-search-value']}>{searchStr}</span>
                            {afterStr}
                        </span>
                    ) : (
                        <span>{strTitle}</span>
                    );
                if (item.children) {
                    return { title, key: item.key, children: loop(item.children) };
                }

                return {
                    title,
                    key: item.key,
                };
            });

        return loop(props.defaultData);
    }, [searchValue, props.defaultData]);

    return (
        <Flex vertical gap="small" style={{ height: '99%' }}>
            {props.showSearch ? (
                <Search  placeholder="Search" onChange={onChange} />
            ) : null}
            <div style={{ height: props.showSearch ? '89%' : '99%' }}>
                <DirectoryTree
                    ref={treeRef}
                    className={styles['ant-tree-directory']}
                    showLine // 显示连接线
                    showIcon={props.showIcon} // 不显示图标
                    defaultExpandAll={true}
                    expandAction="doubleClick"
                    onClick={props.onClick}
                    onDoubleClick={props.onDoubleClick}
                    treeData={treeData}
                    expandedKeys={expandedKeys}
                    onExpand={onExpand}
                    selectedKeys={selectKeys}
                    onSelect={function (selectedKeys) {
                        setSelectKeys(selectedKeys);
                    }}
                    onRightClick={function (info) {
                        setSelectKeys([info.node.key]);
                        if (props.menuOptions) {
                            contextMenu.show({
                                id: menuID,
                                event: info.event,
                                props: info.node.key,
                            });
                        }
                    }}
                    autoExpandParent={autoExpandParent}
                />
                {props.menuOptions ? (
                    <L6Menu menuID={menuID} options={props.menuOptions}></L6Menu>
                ) : null}
            </div>
        </Flex>
    );
});

export default L6Tree;
