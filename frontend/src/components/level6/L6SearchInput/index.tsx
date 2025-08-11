import { AutoComplete } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
interface Option {
    value: string;
}

interface L6SearchInputProps {
    value?: string;
    onChange?: (value: string) => void;
    baseOptions?: string[];
    focusDefaultOpen?: boolean,
    tipOpen?: boolean;
    setTipOpen?: (value: boolean) => void;

    placeholder?: string;
    maxAutoCount?: number;
    onKeyDown?: (event: React.KeyboardEvent<HTMLInputElement>) => void;
    disabled?: boolean,
}
const L6SearchInput: React.FC<L6SearchInputProps> = (props) => {
    const ref = useRef<any>();
    const [filterOptions, setFilterOptions] = useState<Option[]>([]);
    const [open, setOpen] = useState(props.tipOpen ? props.tipOpen : false);

    const wrapOpen = useMemo(() => {
        if (props.setTipOpen) {
            return props.tipOpen;
        }
        return open;
    }, [open, props.tipOpen]);
    const wrapSetOpen = function (v: boolean) {
        if (props.setTipOpen) {
            props.setTipOpen(v);
        } else {
            setOpen(v);
        }
    };

    const handleKeyDown = function (event: React.KeyboardEvent<HTMLInputElement>) {
        if (event.code === 'Escape' || (!event.ctrlKey && event.code === 'Enter')) {
            wrapSetOpen(!wrapOpen);
            event.preventDefault();
            if (ref.current) {
                setTimeout(function () {
                    if (ref.current) {
                        ref.current.focus();
                    }
                }, 100);
            }
        } else {
            if (props.onKeyDown) {
                props.onKeyDown(event);
            }
        }
    };
    const options: Option[] = useMemo(
        () => (props.baseOptions ? props.baseOptions.map((item) => ({ value: item })) : []),
        [props.baseOptions],
    );

    const handleSearch = (val: string) => {
        const res: Option[] = [];
        let count = 0;

        for (let i = options.length - 1; i >= 0; i--) {
            const option = options[i];
            if (option.value.toUpperCase().indexOf(val.toUpperCase()) >= 0) {
                res.push(option);
                count++;
            }
            if (props.maxAutoCount && count > props.maxAutoCount) {
                break;
            }
        }
        setFilterOptions(res);
    };
    const handleBlur = function () {
        wrapSetOpen(false);
    };
    useEffect(() => handleSearch(''), [options]);

    return (
        <AutoComplete
            style={{ width: '100%' }}
            ref={ref}
            autoFocus={true}
            value={props.value}
            options={filterOptions}
            defaultActiveFirstOption={false}
            defaultOpen={wrapOpen}
            open={wrapOpen}
            onFocus={function () {
                if (props.focusDefaultOpen) {
                    wrapSetOpen(props.focusDefaultOpen)
                }
            }}
            backfill={true}
            placeholder={props.placeholder}
            onSearch={handleSearch}
            onChange={props.onChange}
            onBlur={handleBlur}
            onInputKeyDown={handleKeyDown}
            disabled={props.disabled}
        />
    );
};
L6SearchInput.defaultProps = {
    maxAutoCount: 30,
    placeholder: '',
};

export default L6SearchInput;
