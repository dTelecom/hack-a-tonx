import React from 'react';
import {ArrowLeftIcon, Logo} from '../../assets';
import * as styles from './Header.module.scss';
import {ConnectWalletButton} from '../ConnectWalletButton/ConnectWalletButton';
import {observer} from 'mobx-react';
import {appStore} from '../../stores/appStore';
import ProfileBadge from '../ProfileBadge/ProfileBadge';
import classNames from 'classnames';

export const Header = observer(({onBack, title}) => {
  const {currentUser} = appStore;

  const signIn = () => {
    // TODO: sign in
  };

  const signOut = () => {
    // TODO: sign out
    appStore.setCurrentUser(undefined);
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
