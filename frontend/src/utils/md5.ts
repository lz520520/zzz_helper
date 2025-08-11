import { stringToBytes } from '@/utils/byte';
import { Md5 } from 'ts-md5';

export function md5sum(input: string | number[]): string {
    if (typeof input === 'string') {
        input = stringToBytes(input);
    }
    const md5 = new Md5();
    md5.appendByteArray(new Uint8Array(input));

    return String(md5.end());
}
