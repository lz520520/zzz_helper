
import type { IToolPanelComp, IToolPanelParams } from 'ag-grid-community';
import L6Tree from "@/components/level6/L6Tree";
import { createRoot } from "react-dom/client";

const AllTag = "all"
const NotGroupTag = "not_group"

export interface GroupTreeNode {
    title: string,
    key: string,
    count: number,
}
export interface GroupToolPanelParams extends IToolPanelParams {
    groupKey: string;
    treeData: GroupTreeNode[],
}

export class GroupToolPanel implements IToolPanelComp {
    eGui!: HTMLDivElement;
    init(params: GroupToolPanelParams) {
        if (params.groupKey && params.treeData) {
            this.eGui = document.createElement('div');
            this.eGui.style.width  = '100%';
            const root = createRoot(this.eGui);
            root.render(this.calculateStats(params))

            // calculate stats when new rows loaded, i.e. onModelUpdated
            // const renderStats = () => {
            //     root.render(this.calculateStats(params))
            // };
            // params.api.addEventListener('modelUpdated', renderStats);
        }
    }

    getGui() {
        return this.eGui;
    }

    refresh(): void {}

    calculateStats(params: GroupToolPanelParams) {
        // let numGold = 0,
        //     numSilver = 0,
        //     numBronze = 0;
        let treeData = params.treeData;
        for (let i = 0; i < treeData.length; i++) {
            treeData[i].count = 0;
        }
        let count = 0;
        let notGroupCount = 0;
        params.api.forEachNode(function (rowNode) {
            if (!rowNode.data) {
                return
            }
            const value = rowNode.data[params.groupKey];
            if (value !== null) {
                for (let i = 0; i < treeData.length; i++) {
                    if (treeData[i].key == value) {
                        treeData[i].count += 1;
                        break;
                    }
                }
                if (value === "") {
                    notGroupCount += 1;
                }
            }
            count += 1;
        });
        if (!treeData.some(data => data.key == AllTag)) {
            treeData.unshift({
                key: AllTag,
                title: "全部",
                count: count,
            })
            treeData.push({
                key: NotGroupTag,
                title: "未分组",
                count: notGroupCount,
            })
        } else {
            treeData[0].count = count;
            treeData[treeData.length-1].count = notGroupCount;
        }


        return <div style={{ marginRight: '5px', height: '100%' }}>
            <L6Tree defaultData={treeData.map((data) => (
                {
                    title: `${data.title} (${data.count})`,
                    key: data.key,
                }
            ))} onDoubleClick={(e, node)=> {
                let key = node.key;
                let match = "equals"
                if (node.key === AllTag) {
                    key = ""
                } else if (node.key == NotGroupTag) {
                    match = "blank"
                }
                const currentModel =  params.api.getFilterModel();
                currentModel[params.groupKey] = {
                    filter: key,
                    filterType: "text",
                    type: match
                }

                params.api.setFilterModel(currentModel)
                // node.key
            }}></L6Tree>
        </div>;
    }
}
