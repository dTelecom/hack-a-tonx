import React from 'react';
import {ArrowLeftIcon, Logo} from '../../assets';
import * as styles from './Header.module.scss';
import {ConnectWalletButton} from '../ConnectWalletButton/ConnectWalletButton';
import {observer} from 'mobx-react';
import {appStore} from '../../stores/appStore';
import ProfileBadge from '../ProfileBadge/ProfileBadge';
import classNames from 'classnames';
import { connector } from '../../connector';

export const Header = observer(({onBack, title}) => {
  const {currentUser} = appStore;

  const signIn = async () => {
    const walletsList = await connector.getWallets();

    const tonkeeperConnectionSource = {
        universalLink: walletsList[0].universalLink,
        bridgeUrl: walletsList[0].bridgeUrl,
    };

    const universalLink = connector.connect(tonkeeperConnectionSource);
    console.log(universalLink);
  };

  const signOut = () => {
    connector.disconnect();
    // appStore.setCurrentUser(undefined);
  };

  return <div className={classNames(styles.container, title && styles.containerWithTitle)}>
    <div className={styles.logoContainer}>
      <img
        src={Logo}
        alt={'Meet logo'}
      />
    </div>

    {title && (
      <div className={classNames(styles.backContainer, styles.backContainerDesktop)}>
        <button className={styles.backButton}>
          <img
            onClick={onBack}
            src={ArrowLeftIcon}
            alt={'back button'}
          />
        </button>

        <p className={styles.backTitle}>{title}</p>
      </div>
    )}

    <div className={styles.controlContainer}>
      {currentUser ? (
        <ProfileBadge signOut={signOut}/>
      ) : (
        <ConnectWalletButton onClick={signIn}/>
      )}
    </div>

    {title && (
      <div className={classNames(styles.backContainer, styles.backContainerMobile)}>
        <button className={styles.backButton}>
          <img
            onClick={onBack}
            src={ArrowLeftIcon}
            alt={'back button'}
          />
        </button>

        <p className={styles.backTitle}>{title}</p>
      </div>
    )}
  </div>;
});
