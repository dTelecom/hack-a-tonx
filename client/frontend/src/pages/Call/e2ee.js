const worker = new Worker(new URL('../../e2ee_worker.js', import.meta.url), {name: 'E2EE worker'});

export function setKey(rawKey) {
    crypto.subtle.importKey('raw', rawKey, 'AES-GCM', false, ['encrypt','decrypt']).then(key => {
        worker.postMessage({
            operation: 'setKey',
            key,
        });
    });
}

export function setupSenderTransform(sender) {
    if (window.RTCRtpScriptTransform) {
        sender.transform = new window.RTCRtpScriptTransform(worker, {operation: 'encode'});
        return;
    }

    const senderStreams = sender.createEncodedStreams();
    // Instead of creating the transform stream here, we do a postMessage to the worker. The first
    // argument is an object defined by us, the second is a list of variables that will be transferred to
    // the worker. See
    //   https://developer.mozilla.org/en-US/docs/Web/API/Worker/postMessage
    const {readable, writable} = senderStreams;
    worker.postMessage({
        operation: 'encode',
        readable,
        writable,
    }, [readable, writable]);
}

export function setupReceiverTransform(receiver) {
    if (window.RTCRtpScriptTransform) {
        receiver.transform = new window.RTCRtpScriptTransform(worker, {operation: 'decode'});
        return;
    }

    const receiverStreams = receiver.createEncodedStreams();
    const {readable, writable} = receiverStreams;
    worker.postMessage({
        operation: 'decode',
        readable,
        writable,
    }, [readable, writable]);
}
