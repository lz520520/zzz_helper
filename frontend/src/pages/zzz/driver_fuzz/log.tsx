import React, {useEffect, useImperativeHandle, useRef, useState} from "react";
import {EventsOff, EventsOn} from "../../../../wailsjs/runtime";
import {Flex} from "antd";
import globalStyles from "@/utils/style.module.less";
import L6CodeEditor, {L6CodeEditorRef} from "@/components/level6/L6CodeEditor";


const MAX_LOG_SIZE = 1024 * 1024

export type DriverLogRef = {
    setLog:  (msg: string) => void;
};

interface DriverLogProps {
    eventID: string
}

const DriverLogPage = React.forwardRef<DriverLogRef, DriverLogProps>((props, ref) => {
    const [log, setLog] = useState("")

    const logPrint = function (msg: string) {
        setLog(msg);
    };

    useImperativeHandle(ref, () => ({
        setLog: setLog,
    }))

    useEffect(() => {
        EventsOn(`driver_fuzz_${props.eventID}`, (e: string) => {
            logPrint(e)
        })
        return () => {
            EventsOff(`driver_fuzz_${props.eventID}`)
        }
    }, []);
    return (<Flex vertical style={{height: '100%'}} flex={"0 0 140px"}>
        <span className={globalStyles['l6-label']}>日志</span>
        <L6CodeEditor value={log} onChange={setLog} height={'calc( 100% - 32px )'} simpleMode/>
    </Flex>)
})

export default DriverLogPage;
