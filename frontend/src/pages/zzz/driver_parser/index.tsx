import PageWrapTabs from "@/pages/components/PageWrapTabs";
import {Button, Flex, Image, message, TreeDataNode} from "antd";
import Dragger from "antd/es/upload/Dragger";
import {InboxOutlined} from "@ant-design/icons";
import L6CodeEditor from "@/components/level6/L6CodeEditor";
import {Allotment} from "allotment";
import {useEffect, useMemo, useState} from "react";
import { DriverParser, ReadDriverCache} from "../../../../wailsjs/go/main_control/Control";
import {uuid} from "@/utils/uuid";
import {EventsOff, EventsOn} from "../../../../wailsjs/runtime";
import {RcFile, UploadFile} from "antd/es/upload";
import { zzz_models} from "../../../../wailsjs/go/models";
import {md5sum} from "@/utils/md5";
import {Base64} from "js-base64";
import DBTable from "@/pages/components/DBTable/DBTable";
import globalStyles from "@/utils/style.module.less";

interface FileInfo {
    id: string,
    file: UploadFile,
}

const Content: React.FC<{id?: string}> = ({id}) => {
    const [fileInfos, setFileInfos] = useState<FileInfo[]>([]);
    const eventID = useMemo(() => uuid(),[])
    const [log, setLog] = useState("")
    const [imageBas64, setImageBase64] = useState("")

    const [previewOpen, setPreviewOpen] = useState(false);
    const [previewImage, setPreviewImage] = useState('');
    const [ocr, setOcr] = useState("")

    const logPrint = function (msg: string) {
        setLog((v) => v + '\r\n' + msg);
    };


    const parserHandle = async () => {
        if (!fileInfos.length) {
            message.warning("图片列表为空")
            return
        }
        const files: Array<zzz_models.FileInfo> = []
        for (let info of fileInfos) {
            const buffer = await (info.file as RcFile).arrayBuffer();
            const data = Array.from(new Uint8Array(buffer))
            const id = md5sum(data)
            files.push({
                id,
                data
            })
            info.id = id
        }
        const resp = await DriverParser(eventID, files, ocr)
        if (!resp.status) {
            message.error(resp.err)
            return
        }
        message.success("解析完毕")
        const idSet = new Set(resp.ids);
        setFileInfos(prev => prev.filter(item => !idSet.has(item.id)));
        setOcr("")
    }
    const diskTreeDoubleClickHandle = async(currentRow: Record<string, any>) => {
        const resp = await  ReadDriverCache(`${currentRow.id}`)

        if (!resp.status) {
            message.error('获取失败, ' + resp.err);
            return;
        }
        setImageBase64(resp.bytes.toString())
    }

    useEffect(() => {
        EventsOn(`driver_parser_${eventID}`, (e: string) => {
            logPrint(e)
        })
        return () => {
            EventsOff(`driver_parser_${eventID}`)
        }
    }, []);



    return (<Allotment>
        <Allotment.Pane preferredSize="35%">
            <Flex vertical gap={"small"} style={{height: '100%'}}>

                <div style={{height: '60%', width: '100%'}}>
                    <DBTable moduleName={"driver_cache"} keyName={"id"} operation={{
                        MultiDelete: true,
                        Detail: true,
                    }} onRowClick={diskTreeDoubleClickHandle}/>

                </div>
                <div className={globalStyles['l6-image']} style={{height: '40%'}}>
                    <Image
                        src={`data:image/png;base64,${imageBas64}`}
                        fallback="data:image/png;base64,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"
                    />
                </div>
            </Flex>


        </Allotment.Pane>
        <Allotment.Pane preferredSize="25%">
            <Flex gap={"small"} vertical style={{ height: "100%", overflow: 'auto' }}>
                <Dragger
                    name={"file"}
                    fileList={fileInfos.map((item) => item.file)}
                    multiple={true}
                    height={200}
                    accept={'image/*'}
                    beforeUpload={async (file) => {
                        // 阻止自动上传
                        const buffer = await file.arrayBuffer();

                        const newFile = file as UploadFile
                        newFile.thumbUrl = `data:image/png;base64,${Base64.fromUint8Array(new Uint8Array(buffer))}`
                        setFileInfos((prevList) => [...prevList, {id: "", file: newFile}]);
                        return false;
                    }}
                    onRemove={(file) => {
                        setFileInfos((prevList) => prevList.filter((item) => item.file.uid !== file.uid));
                    }}
                    customRequest={function ({onSuccess}) {
                        if (onSuccess) onSuccess("ok");
                    }}
                    onPreview={(file) => {
                        setPreviewImage(file.thumbUrl as string);
                        setPreviewOpen(true);
                    }}

                    pastable={true}
                    listType={'picture'}
                >
                    <p className="ant-upload-drag-icon">
                        <InboxOutlined />
                    </p>
                    <p className="ant-upload-text">点击或者拖拽图片</p>
                </Dragger>
                {previewImage && (
                    <Image
                        wrapperStyle={{ display: 'none' }}
                        preview={{
                            visible: previewOpen,
                            onVisibleChange: (visible) => setPreviewOpen(visible),
                            afterOpenChange: (visible) => !visible && setPreviewImage(''),
                        }}
                        src={previewImage}
                    />
                )}
                <Button
                    type={"primary"}
                    onClick={parserHandle}
                    disabled={fileInfos.length === 0}
                >提交</Button>
                <Button
                    danger
                    onClick={function () {
                        setFileInfos([])
                    }}
                >清空</Button>
                <L6CodeEditor value={ocr} onChange={(item) => setOcr(item)}  simpleMode/>
            </Flex>

        </Allotment.Pane>
        <L6CodeEditor style={{width: '100%', margin: '10px'}} value={log} onChange={(code) => setLog(code)} scrollMode/>

    </Allotment>)
}

const DriverParserPage: React.FC =() => {
    return (
        <PageWrapTabs
            tag={"driver_parser"}
            content={(id?: string)=> <Content id={id}/>}
        />
    );
}


export default DriverParserPage;
