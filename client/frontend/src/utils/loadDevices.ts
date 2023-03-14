export const loadDevices = (callback: (devices: MediaDeviceInfo[]) => void) => {
  navigator.mediaDevices.enumerateDevices()
    .then((devices) => {
      callback(devices)
    })
    .catch((err) => {
      console.error(`${err.name}: ${err.message}`);
    });
}