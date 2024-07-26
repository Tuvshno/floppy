const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electron', {
    logFilePath: (filePath) => ipcRenderer.send('log-file-path', filePath),
    executeBinary: (filePath) => ipcRenderer.send('execute-binary', filePath),
});
