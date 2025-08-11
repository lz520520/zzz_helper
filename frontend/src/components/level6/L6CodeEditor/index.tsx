import { sleep } from '@/utils/time';
import { Extension } from '@codemirror/state';
import * as alls from '@uiw/codemirror-themes-all';
import CodeMirror, { BasicSetupOptions, ViewUpdate } from '@uiw/react-codemirror';
import React, {CSSProperties, useEffect, useImperativeHandle, useMemo, useRef, useState} from 'react';
import { useImmer } from 'use-immer';
import { langNames, loadLanguage } from './langs';
import './styles.less';
import styles from './styles.module.less';
import { useTheme } from './useTheme';
import {Flex} from "antd";
import {uuid} from "@/utils/uuid";
import {ReactCodeMirrorRef} from "@uiw/react-codemirror/src";
import type {EditorView} from "@codemirror/view";

const themeOptions = ['dark', 'light']
    .concat(Object.keys(alls))
    .filter((item) => typeof alls[item as keyof typeof alls] !== 'function')
    .filter((item) => !/^(defaultSettings)/.test(item as keyof typeof alls));

interface CodeState {
    lines: number;
    cursor: number;
    selected: number;
    length: number;
}

interface L6CodeEditorProps {
    value?: string;
    onChange?: (code: string) => void;
    language?: string;
    simpleMode?: boolean;
    height?: string;
    scrollMode?: boolean; // 滚动模式，会自动滚动到最新的消息
    style?: CSSProperties;
    init?: () => void,
    placeholder?: string;
}
export interface L6CodeEditorRef  {
    appendText: (text: string) => void;
    setText: (text: string) => void;
    getTextLength: () => number,
}

const L6CodeEditor = React.forwardRef<L6CodeEditorRef, L6CodeEditorProps>((props, ref) => {

    const editorRef = useRef<ReactCodeMirrorRef>(null);
    const [mode, setMode] = useState(props.language ? props.language : 'l6_event');
    const [placeholder, setPlaceholder] = useState( props.placeholder? props.placeholder: 'Please enter the code.');
    const [autofocus, setAutofocus] = useState(false);
    const [editable, setEditable] = useState(true);
    const [scroll, setScroll] = useState(!!props.scrollMode);
    const id = useMemo(() => uuid(), [])

    const { theme, setTheme } = useTheme();
    const [extensions, setExtensions] = useState<Extension[]>();
    const [basicSetup, setBasicSetup] = useState<BasicSetupOptions>({
        crosshairCursor: false,
        tabSize: 4,
    });
    const [viewState, updateViewState] = useImmer<CodeState>({
        lines: 0,
        cursor: 0,
        selected: 0,
        length: 0,
    });

    const appendText = (text: string) => {
        const view = editorRef.current?.view as EditorView;
        view.dispatch({
            changes: { from: view.state.doc.length, insert:  text},
        });
    };

    const setText = (text: string) => {
        const view = editorRef.current?.view as EditorView;
        if (!view) return;

        view.dispatch({
            changes: {
                from: 0,
                to: view.state.doc.length,
                insert: text,
            },
        });
    };
    const getTextLength = () => {
        const view = editorRef.current?.view as EditorView;
        if (!view) return 0;
        return view.state.doc.length;
    };
    useImperativeHandle(ref, () => {
        return {
            appendText: appendText,
            setText: setText,
            getTextLength: getTextLength,
        }
    }, []);

    useEffect(() => {
        if (editorRef.current && editorRef.current.view) {
            if (props.init) {
                props.init()
            }
        }
    }, [editorRef.current]);

    const scrollFunc = function () {
        const container = document.getElementById(id);
        if (!container) {
            return;
        }
        const items = container.getElementsByClassName('cm-scroller');
        for (let i = 0; i < items.length; i++) {
            const item = items.item(i) as HTMLDivElement;
            if (item.scrollHeight > 0) {
                item.scrollTop = item.scrollHeight;

                break;
            }
        }
    };




    const handleChange = function (value: string, viewUpdate: ViewUpdate): boolean {
        if (scroll) {
            const tmpView = viewUpdate as any;
            // flags 6表示有输入，transactions表示无事务处理，可能是滚动
            if (
                (tmpView['flags'] && tmpView['flags'] != 6) ||
                viewUpdate.transactions.length == 0
            ) {
                return true;
            }
            scrollFunc();
        }

        return true;
    };

    useEffect(() => {
        const lang = loadLanguage(mode);
        if (lang) {
            setExtensions(lang);
        }
    }, [mode]);
    useEffect(() => {
        const fetchData = async () => {
            if (scroll) {
                await sleep(500);
                // 单独设置是因为初始化时获取不到cm-scroller，需要等待，这里异步执行，可以等初始化完后再滚动
                scrollFunc();
            }
        };
        fetchData();
    }, [scroll]);


    return (
        <div
            className="wmde-markdown-var"
            style={{ ...props.style, height: props.height ? props.height : '100%', fontSize: 12 }}
        >
            {!props.simpleMode? <>
                <div className={styles['toolbar']}>
                    <div className="item">
                        <label htmlFor="language">language:</label>
                        <select
                            className={'select'}
                            id="language"
                            value={mode}
                            onChange={(e) => setMode(e.target.value)}
                        >
                            {langNames.map((key) => (
                                <option value={key} key={key}>
                                    {key}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div className="item">
                        <Flex gap={'small'}>
                            <div className="item">
                                <label htmlFor="scroll">scroll:</label>
                                <input
                                    type="checkbox"
                                    id="disabled"
                                    checked={scroll}
                                    onChange={(evn) => setScroll(evn.target.checked)}
                                />
                            </div>
                            <div className="item">
                                <label htmlFor="editable">editable:</label>
                                <input
                                    type="checkbox"
                                    id="disabled"
                                    checked={editable}
                                    onChange={(evn) => setEditable(evn.target.checked)}
                                />
                            </div>
                        </Flex>

                    </div>

                    {/*<div className={styles["item"]}>*/}
                    {/*  <label htmlFor="theme">theme:</label>*/}
                    {/*  <select*/}
                    {/*    id="theme"*/}
                    {/*    className={styles["select"]}*/}
                    {/*    value={theme as string}*/}
                    {/*    onChange={(evn) => {*/}
                    {/*      setTheme(evn.target.value as ReactCodeMirrorProps['theme']);*/}
                    {/*    }}*/}
                    {/*  >*/}
                    {/*    {themeOptions.map(key => <option value={key} key={key}>{key}</option>)}*/}
                    {/*  </select>*/}
                    {/*</div>*/}
                </div>
                <div className={styles['divider']}></div>
            </> : null}
            <CodeMirror
                ref={editorRef}
                value={props.value}
                id={id}
                style={{height: props.simpleMode ? "100%" : 'calc( 100% - 4rem - 2px )'}}
                height={`100% !important`}
                className={styles['editor']}
                // @ts-ignore
                theme={alls[theme as keyof typeof alls] || theme}
                editable={editable}
                extensions={extensions}
                autoFocus={autofocus}
                basicSetup={basicSetup}
                placeholder={placeholder}
                onUpdate={(viewUpdate) => {
                    const update = viewUpdate as Record<any, any>;
                    if ((update['flags'] && update['flags'] === 6)) {
                        handleChange('', viewUpdate);
                    }
                    updateViewState((draft) => {
                        // selected
                        const ranges = viewUpdate.state.selection.ranges;
                        draft.selected = ranges.reduce(
                            (plus, range) => plus + range.to - range.from,
                            0,
                        );
                        draft.cursor = ranges[0].anchor;
                        // length
                        draft.length = viewUpdate.state.doc.length;
                        draft.lines = viewUpdate.state.doc.lines;
                    });
                }}
                onChange={(value: string, viewUpdate: ViewUpdate) => {
                    handleChange(value, viewUpdate);
                    if (props.onChange) {
                        props.onChange(value);
                    }
                    // https://github.com/uiwjs/react-codemirror/issues/449
                    // setCode(val)
                }}
            />
            {!props.simpleMode? <>
                <div className={styles['divider']}></div>
                <div className={styles['footer']}>
                    <div className="infos">
                        <span className="item">Length: {viewState.length}</span>
                        <span className="item">Lines: {viewState.lines}</span>
                        <span className="item">Cursor: {viewState.cursor}</span>
                        <span className="item">Selected: {viewState.selected}</span>
                    </div>
                </div>
            </> : null
            }
        </div>
    );
});
export default L6CodeEditor;
