import React, {ChangeEventHandler, useEffect, useMemo, useState} from "react";
import {Checkbox, Flex, Input} from "antd";
import {ValueType} from "rc-input/lib/interface";


interface L6CheckInputProps  {
    onChange?: ChangeEventHandler<HTMLInputElement> | undefined;
    width?: number | string | undefined;
    value?: ValueType;
    vertical?: boolean
}
const L6CheckInput: React.FC<L6CheckInputProps> = (props) => {
    const [check, setCheck] = useState(!!props.value)
    const onChange: ChangeEventHandler<HTMLInputElement> = (e) => {
        if (!props.onChange){
            return
        }
        if (!check) {
            e.target.value = ""
        }
        props.onChange(e)
    }

    useEffect(() => {
        if (!props.onChange){
            return
        }
        const e: React.ChangeEvent<HTMLInputElement> = {
            target: {
                value: check?props.value: "",
            }
        } as React.ChangeEvent<HTMLInputElement>;
        props.onChange(e)
    }, [check]);


    return  <Flex gap="small" vertical={props.vertical}>
        <Checkbox checked={check} onChange={(e) => {
            setCheck(e.target.checked)
        }}/>
        <Input  value={props.value}  onChange={onChange} width={props.width}/>
    </Flex>
}

export default L6CheckInput;
