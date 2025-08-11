export type CommonOperationFunc = (data: Record<string, any>) => Promise<void>;
export type SingleOperationFunc = (condition: Record<string, any>) => Promise<void>;
export type MultiOperationFunc = (datas: Record<string, any>[]) => Promise<void>;

type ListOperationFunc = (tips?: boolean) => Promise<void>;

export interface OperationSupport {
    List: ListOperationFunc;
    Create?: false | CommonOperationFunc;
    Update?: false | CommonOperationFunc;

    Query?: false | CommonOperationFunc;
    Delete?: false | CommonOperationFunc;
    MultiDelete?: false | MultiOperationFunc;
    Clear?: false | SingleOperationFunc;
    Detail?: boolean;
    CopyColumns?: boolean,
    Exports?: string
}


export interface ExtractInfo {
    name: string,
    expression: string,
    post?:(result: string) => string
}



export interface ExpressionInfo {
    type: string,
    field?: string,
    value?: string,
}
