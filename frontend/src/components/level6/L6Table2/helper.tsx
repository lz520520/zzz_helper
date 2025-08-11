import type {CellStyleFunc, ColDef,  SortDirection} from 'ag-grid-community';
import { CommonConfigOption } from '@/models/interface';
import {ExpressionInfo} from "./types";

export const getColumnTypeFromConfig = (
    options: CommonConfigOption[],
    customCell?: Record<string, CellStyleFunc> | CellStyleFunc,
    noellipsis = false,
    customHide?: Record<string, boolean>
) => {
    const columns: ColDef[] = [];
    let colIndex = 0;
    const renderCell = function (key: string) {
        if (!customCell) {
            return noellipsis ? { whiteSpace: 'pre' } : undefined;
        }
        if (typeof customCell === 'function') {
            return customCell as CellStyleFunc;
        } else if (key in customCell) {
            return customCell[key];
        }
        return noellipsis ? { whiteSpace: 'pre' } : undefined;
    };
    options.map((option) => {
        if (option.show_in_table) {
            let col: ColDef = {
                field: option.key,
                headerName: option.title,
                width: option.col_width ? option.col_width : 150,
                pinned: !!option.col_fixed,
                initialHide: customHide && option.key in customHide? customHide[option.key] : option.col_hidden,
                autoHeight: noellipsis,
                filter: option.custom_filter_dropdown ? 'agTextColumnFilter' : undefined,
                filterParams: {
                    buttons: ['apply', 'reset'],
                    closeOnApply: true,
                    maxNumConditions: 5,
                },
                editable: option.edit,
                sort: option.sort? option.sort as SortDirection: null,
                tooltipField: option.key,
                cellStyle: renderCell(option.key),
                rowGroup: !!option.row_group,
            };

            columns.push(col);
        }
    });
    return columns;
};




// 预编译表达式
export const compileExpression = (expression: string): ExpressionInfo[] => {
    const regex = /\{\{(.*?)\}\}/g;
    const parts:ExpressionInfo[] = [];
    let match;
    let lastIndex = 0;

    while ((match = regex.exec(expression)) !== null) {
        // 处理 {{ 之前的文本
        if (match.index > lastIndex) {
            parts.push({ type: 'text', value: expression.slice(lastIndex, match.index) });
        }

        // 处理 {{ }} 内的条件或变量
        const expPart = match[1];
        if (expPart.includes('=')) {
            const [field, value] = expPart.split('=');
            parts.push({ type: 'condition', field, value: value.replace(/['"]/g, '') });
        } else {
            parts.push({ type: 'variable', field: expPart });
        }

        lastIndex = regex.lastIndex;
    }

    // 添加剩余的文本
    if (lastIndex < expression.length) {
        parts.push({ type: 'text', value: expression.slice(lastIndex) });
    }

    return parts;
};


// 根据预编译后的表达式，处理每一行的数据
export const evaluateExpression = (compiledExpression: ExpressionInfo[], row: Record<string, any>):string => {
    let result = '';

    for (const part of compiledExpression) {
        if (part.type === 'text') {
            result += part.value; // 直接拼接文本
        } else if (part.type === 'variable' && part.field) {
            result += row[part.field] || ''; // 替换变量
        } else if (part.type === 'condition' && part.field) {
            // 条件不满足则返回空字符串
            if (!row[part.field].toLowerCase().includes(part.value)) {
                return '';
            }
            result += row[part.field] || ''; // 替换变量
        }
    }

    return result;
};
