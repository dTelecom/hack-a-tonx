import React from 'react';
import {Flex} from '@chakra-ui/react';
import {dTelecomLogo} from '../../assets';
import styles from './Footer.module.scss';

const Footer = () => {
  return (
    <Flex
      className={styles.container}
      justifyContent={'center'}
    >
      <p className={styles.text}>Powered by</p><a
      href="https://dtelecom.org/"
      rel="noreferrer"
      target="_blank"
    ><img
      src={dTelecomLogo}
      alt="dTelecom logo"
    /></a>
    </Flex>
  );
};

export default Footer;