import {Editor} from '@monaco-editor/react';
import {useEffect, useState, useRef} from 'react';

const CodeEditor = () => {
    const [code, setCode] = useState('');
    const ws = useRef(null);

    useEffect(() => {
        ws.current = new WebSocket('ws://localhost:8080/ws');

        ws.current.onopen = () => {
            console.log('Connected to server');
        };

        window.addEventListener("focus", () => {
            fetch("http://localhost:8080/code")
                .then(response => response.text())
                .then(data => setCode(data));
        })

        ws.current.onmessage = (event) => {
            setCode(event.data);
        };

        ws.current.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        ws.current.onclose = () => {
            console.log('Disconnected from server');
        };

    }, []);

    const handleCodeChange = (value) => {
        setCode(value);
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(value);
        } else {
            console.error('WebSocket not connected');
        }
    };

    return (
        <div>
            <Editor
                height="95vh"
                defaultLanguage="javascript"
                defaultValue={code}
                theme="vs-dark"
                value={code}
                onChange={handleCodeChange}
            />
        </div>
    );
}

export default CodeEditor;