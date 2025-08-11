export async function readChunk(file: Blob, offset: number, length: number): Promise<Uint8Array> {
    return new Promise<Uint8Array>((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (event) => {
            const array = new Uint8Array(reader.result as ArrayBuffer);
            resolve(array);
        };
        reader.onerror = reject;

        const blobSlice = file.slice(offset, offset + length);
        reader.readAsArrayBuffer(blobSlice);
    });
}

export function formatBytes(bytes: number) {
    if (bytes === 0) return '0 B';

    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));
    const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));

    return `${size} ${sizes[i]}`;
}

export function stringToBytes(str: string): number[] {
    const bytes: number[] = [];

    for (let i = 0; i < str.length; i++) {
        // 获取字符串中每个字符的 ASCII 值并添加到数组中
        bytes.push(str.charCodeAt(i));
    }

    return bytes;
}
