const dropZone = document.getElementById('drop-zone');

dropZone.addEventListener('dragover', (event) => {
    event.preventDefault();
    dropZone.classList.add('drag-over');
});

dropZone.addEventListener('dragleave', () => {
    dropZone.classList.remove('drag-over');
});

dropZone.addEventListener('drop', (event) => {
    event.preventDefault();
    dropZone.classList.remove('drag-over');

    const files = event.dataTransfer.files;
    for (const file of files) {
        const filePath = file.path;
        window.electron.logFilePath(filePath);
        window.electron.executeBinary(filePath);
    }
});
