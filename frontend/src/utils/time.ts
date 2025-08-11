const minute = 60;
const hour = 60 * minute;
const day = hour * 24;

export function sleep(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
export function TimeSimpleFormat(t: number): string {
    let timeStr = '';
    switch (true) {
        case t < minute:
            timeStr = `${t}s`;
            break;
        case t < hour:
            timeStr = `${Math.floor(t / minute)}m${Math.floor(t % minute)}s`;
            break;
        case t < day:
            timeStr = `${Math.floor(t / hour)}h${Math.floor((t % hour) / minute)}m`;
            break;
        default:
            timeStr = `${Math.floor(t / day)}d${Math.floor((t % day) / hour)}h`;
            break;
    }
    return timeStr;
}

export function TimeFormat(t: number): string {
    let timeStr = '';
    switch (true) {
        case t < minute:
            timeStr = `${t}s`;
            break;
        case t < hour:
            timeStr = `${Math.floor(t / minute)}m${Math.floor(t % minute)}s`;
            break;
        case t < day:
            timeStr = `${Math.floor(t / hour)}h${Math.floor((t % hour) / minute)}m${Math.floor(
                (t % hour) % minute,
            )}s`;
            break;
        default:
            timeStr = `${Math.floor(t / day)}d${Math.floor((t % day) / hour)}h${Math.floor(
                ((t % day) % hour) / minute,
            )}m${Math.floor(((t % day) % hour) % minute)}s`;
            break;
    }
    return timeStr;
}


export function TimestampFormat(timestamp: string): string {
  const date = new Date(timestamp);

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");

  const formattedTimestamp = `${year}-${month}-${d} ${hours}:${minutes}:${seconds}`;
  return formattedTimestamp
}



export function getFormattedTimestamp(): string {
    const date = new Date();

    // 获取各个部分并格式化为两位数
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}
