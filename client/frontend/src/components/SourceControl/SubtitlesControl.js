import React, {useCallback} from 'react';
import {ArrowDisabledDown, ArrowDownIcon, EnabledTickIcon, SubOffIcon, SubOnIcon} from '../../assets';
import {IconButton, Popover, PopoverBody, PopoverContent, PopoverTrigger, useDisclosure} from '@chakra-ui/react';
import styles from './SourceControl.module.scss';
import classNames from 'classnames';

const SubtitlesControl = ({selected, devices, enabled, onChange, toggleMute, isCall}) => {
  const enabledIcon = SubOnIcon;
  const disabledIcon = SubOffIcon;
  const enableText = 'Enable subtitles';
  const disableText = 'Disable subtitles';
  const icon = SubOnIcon;
  const {isOpen, onClose, onOpen} = useDisclosure();

  const onMuteClick = useCallback((e) => {
    e.stopPropagation();
    toggleMute();
  }, [toggleMute]);

  const text = enabled ? disableText : enableText;

  return (
    <Popover
      closeOnBlur={false}
      isOpen={isOpen}
      onOpen={onOpen}
      onClose={onClose}
      placement="top-start"
    >
      <PopoverTrigger>
        <div className={classNames(styles.container, !enabled && styles.containerDisabled, isCall && styles.isCall)}>
          <IconButton
            className={styles.iconButton}
            onClick={onMuteClick}
            aria-label={`${text} button`}
          >
            <img
              src={enabled ? enabledIcon : disabledIcon}
              alt={`${text} icon`}
            />
          </IconButton>

          <div className={styles.divider}/>

          <IconButton
            className={styles.iconButton}
            onClick={onOpen}
            aria-label={'select source'}
          >
            <img
              className={classNames(isOpen && styles.arrowOpen)}
              src={enabled ? ArrowDownIcon : ArrowDisabledDown}
              alt={'select source icon'}
            />
          </IconButton>
        </div>
      </PopoverTrigger>

      <PopoverContent style={{outline: 'none'}}>
        <PopoverBody>
          <div className={styles.popOver}>
            {devices?.length > 0 ? devices.map(device => (
              <button
                onClick={() => {
                  onChange(device.deviceId);
                  onClose();
                }}
                key={device.deviceId}
                className={styles.popOverItem}
              >
                <img
                  alt={'source icon'}
                  src={selected === device.deviceId ? EnabledTickIcon : icon}
                />
                <p>{device.label}</p>
              </button>
            )) : null}

            <div className={styles.popOverDisable}>
              <button
                onClick={() => {
                  toggleMute();
                  onClose();
                }}
              >{text}</button>
            </div>
          </div>
        </PopoverBody>
      </PopoverContent>
    </Popover>
  );
};

export default SubtitlesControl;
