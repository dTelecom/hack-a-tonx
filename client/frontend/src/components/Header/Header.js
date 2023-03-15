import React from 'react';
import {Logo, LogoMini} from '../../assets';
import styles from './Header.module.scss';
import {observer} from 'mobx-react';
import classNames from 'classnames';
import {Box, Flex} from '@chakra-ui/react';
import {useBreakpoints} from '../../hooks/useBreakpoints';

export const Header = observer(({children, centered, title, isMiniLogo}) => {
  const {isMobile} = useBreakpoints();

  return (
    <>
      <div className={classNames(styles.container, centered && styles.containerCentered)}>
        <Flex
          flexGrow={1}
          width={!isMobile ? '25%' : 'initial'}
          justifyContent={centered && !children ? 'center' : 'initial'}
        >
          <div className={styles.logoContainer}>
            <img
              src={isMiniLogo ? LogoMini : Logo}
              alt={'dTelecom logo'}
            />
          </div>
        </Flex>

        {!isMobile && (
          <Flex
            flexGrow={1}
            alignItems={'center'}
            justifyContent={'center'}
            width={'50%'}
          >
            {title && (
              <h1 className={styles.title}>
                {title}
              </h1>
            )}
          </Flex>
        )}

        {isMobile && isMiniLogo && title && (
          <Flex className={styles.mobileMiniTitle}>
            <span>{title}</span>
          </Flex>
        )}


        {(!isMobile || (isMobile && children)) && (
          <Flex
            flexGrow={1}
            width={!isMobile ? '25%' : undefined}
          >
            {children && (
              <div className={styles.controlContainer}>
                {children}
              </div>
            )}
          </Flex>
        )}
      </div>
      {isMobile && !isMiniLogo && title && (
        <Box>
          <h3 className={styles.mobileTitle}>{title}</h3>
        </Box>
      )}
    </>
  );
});