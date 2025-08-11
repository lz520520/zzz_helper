import { StreamLanguage, StreamParser } from '@codemirror/language';
import { simpleMode } from '@codemirror/legacy-modes/mode/simple-mode';
// rust
// node_modules/@uiw/codemirror-theme-monokai/src/color.ts
const stream: StreamParser<unknown> = simpleMode({
    start: [
        //
        /*
         * builtin number: yellow
         * atom: orange
         * meta: grey
         * operator: blue
         * labelName： deep blue
         * keyword: pink
         * */

        // 2006-01-02 15:04:05
        { regex: /^\[\d{4}\-\d{2}\-\d{2} \d{2}:\d{2}:\d{2}\]/, token: 'comment' },
        // orange
        { regex: /(\s|^)accelerator\>/, token: 'variable' },

        // green success
        { regex: /(\s|^)\[\+\]/, token: 'className' },
        { regex: /\s\[sent\]/, token: 'heading' },

        // red error
        { regex: /(\s|^)\[\-\]/, token: 'invalid' },
        { regex: /\s\[error\]/, token: 'invalid' },

        // blue info
        { regex: /(\s|^)\[\*\]/, token: 'labelName' },
        { regex: /\s\[prepare\]/, token: 'labelName' },

        // yellow warn
        { regex: /(\s|^)\[\!\]/, token: 'string' },
        // white text
        // {regex: /[\w\[\]\\\u4e00-\u9fa5]+/, token: "cursor"},

        { regex: /\s\[result\]/, token: 'variable' },
    ],
    comment: [
        { regex: /.*?\*\//, token: 'comment', next: 'start' },
        { regex: /.*/, token: 'comment' },
    ],
    languageData: {
        name: 'rust',
        // dontIndentStates: ["comment"],
        // indentOnInput: /^\s*\}$/,
        commentTokens: { line: '//', block: { open: '/*', close: '*/' } },
    },
});
export const l6_event = () => StreamLanguage.define(stream);
