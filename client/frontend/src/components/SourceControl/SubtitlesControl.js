import React, {useCallback} from 'react';
import {SubOffIcon, SubOnIcon} from '../../assets';
import {Flex, IconButton} from '@chakra-ui/react';
import styles from './SourceControl.module.scss';
import classNames from 'classnames';
import {useBreakpoints} from '../../hooks/useBreakpoints';

const SubtitlesControl = ({enabled, toggleMute, isCall}) => {
  const enabledIcon = SubOnIcon;
  const disabledIcon = SubOffIcon;
  const enableText = 'Enable subtitles';
  const disableText = 'Disable subtitles';
  const {isMobile} = useBreakpoints();

  const onMuteClick = useCallback((e) => {
    e.stopPropagation();
    toggleMute();
  }, [toggleMute]);

  const text = enabled ? disableText : enableText;

  return (
    <Flex
      width={isMobile ? 64 : 80}
      justifyContent={'center'}
      className={classNames(styles.container, !enabled && styles.containerDisabled, isCall && styles.isCall)}
    >
      <IconButton
        // className={styles.iconButton}
        onClick={onMuteClick}
        aria-label={`${text} button`}
        style={{
          justifyContent: 'center'
        }}
      >
        <img
          src={enabled ? enabledIcon : disabledIcon}
          alt={`${text} icon`}
        />
      </IconButton>
    </Flex>
  );
};

export default SubtitlesControl;
