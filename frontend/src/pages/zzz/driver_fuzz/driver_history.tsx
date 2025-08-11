import React from "react";
import type {UpdateDataType} from "@/pages/components/DBTable/types";
import {Drawer, Flex, Image} from "antd";
import DBTable from "@/pages/components/DBTable/DBTable";


interface DriverSelectProps {
    open: boolean;
    onChange: (value: boolean) => void
    onRowDoubleClick?: (currentRow: Record<string, any>, update: UpdateDataType) => Promise<void>;
}

const DriverHistoryPage: React.FC<DriverSelectProps> = (props) => {
    return (<Drawer forceRender open={props.open}  onClose={() => props.onChange(false)} width={'50%'}>
        <DBTable moduleName={"driver_fuzz"} keyName={"id"} operation={{
            Detail: true,
        }} onRowDoubleClick={props.onRowDoubleClick}/>
    </Drawer>)
}

export default DriverHistoryPage;

