import React, {useEffect, useState} from "react";
import {useImmer} from "use-immer";
import {PageContainer} from "@ant-design/pro-components";
import {Tabs} from "antd";
import globalStyles from '@/utils/style.module.less';
import {uuid} from "@/utils/uuid";
import {GetTasks} from "../../../wailsjs/go/main_control/Control";
import {common_model} from "../../../wailsjs/go/models";


interface TabItem {
    key: string;
    label: string;
    children: React.ReactNode;
    closable?: boolean;
}

export interface RailPageWrapTabsProps  {
    tag: string
    content: (id?: string) => React.ReactNode;

}
const PageWrapTabs: React.FC<RailPageWrapTabsProps> = (props) => {
    const [tabItems, updateTabItems] = useImmer<TabItem[]>([]);
    const [tabIndex, setTabIndex] = useState(0);
    const [activeTab, setActiveTab] = useState<string>('');
    const getCurrTab = (newActiveTab: string) => tabItems.find((item) => item.key === newActiveTab);

    // 切换 Tab
    const switchTab = (newActiveTab: string) => {
        const currTab = getCurrTab(newActiveTab);
        if (currTab) {
            setActiveTab(newActiveTab);
        }
    };
    // 处理tab增删
    const handleTableEdit = (
        e: React.MouseEvent | React.KeyboardEvent | string,
        action: 'add' | 'remove',
    ) => {
        if (action === 'add') {
            let key = uuid();
            let id = undefined;
            if (typeof e === "string" && e) {
                id = e
            }
            updateTabItems((draft) => {
                draft.push({
                    key: key,
                    label: `${props.tag}-[${tabIndex}]`,
                    children: props.content(id),
                });
            });
            setTabIndex((i) => i + 1);
            setActiveTab(key);
        } else {
            let key = e as string;

            let newActiveTab = activeTab;
            let lastIndex = -1;
            tabItems.forEach((item, i) => {
                if (item.key === key) {
                    lastIndex = i - 1;
                }
            });
            const newPanes = tabItems.filter((item) => item.key !== key);

            if (newPanes.length && newActiveTab === key) {
                if (lastIndex >= 0) {
                    newActiveTab = newPanes[lastIndex].key;
                } else {
                    newActiveTab = newPanes[0].key;
                }
            }
            if (newPanes.length <= 1) {
                setTabIndex(0);
            }

            updateTabItems((draft) => {
                return newPanes;
            });
            setActiveTab(newActiveTab);
        }
    };
    // 首次触发
    useEffect(() => {
        const fetchData = async () => {
            const resp = await GetTasks(new common_model.CommonReq({msg: props.tag}))
            if (resp.status && resp.tasks.length){
                resp.tasks.forEach(item => {
                    handleTableEdit(item.task_id, "add")
                })
            } else {
                handleTableEdit('', 'add');
            }
        }
        fetchData()
    }, []);
    // 不允许关闭最后一个标签
    useEffect(() => {
        if (tabItems.length === 1) {
            updateTabItems((draft) => {
                draft[0].closable = false;
            });
        } else {
            updateTabItems((draft) => {
                for (let i = 0; i < draft.length; i++) {
                    draft[i].closable = true;
                }
            });
        }
    }, [tabItems]);
    return (
        <PageContainer
            header={{ title: '', breadcrumb: {} }}
            className={globalStyles['l6-page-container']}
        >
            <Tabs
                className={globalStyles['l6-tabs']}
                type="editable-card"
                onChange={switchTab}
                activeKey={activeTab}
                onEdit={handleTableEdit}
                items={tabItems}
            ></Tabs>
        </PageContainer>
    );

}

export default PageWrapTabs;
