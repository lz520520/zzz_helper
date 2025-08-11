import { ReactCodeMirrorProps } from '@uiw/react-codemirror';
import { useState } from 'react';

export function useTheme(
    name: ReactCodeMirrorProps['theme'] = 'monokai' as ReactCodeMirrorProps['theme'],
) {
    const [theme, setTheme] = useState<ReactCodeMirrorProps['theme']>(name);
    return { theme, setTheme };
}
