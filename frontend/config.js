function getConfig(env) {
  switch (env) {
    case 'mainnet':
      return {};
    case 'local':
      return {};
    default:
      throw Error(`Unconfigured environment '${env}'. Can be configured in src/config.js.`);
  }
}

module.exports = getConfig;
