import React, {useRef} from 'react';
import styles from './Input.module.scss';
import {Box} from '@chakra-ui/react';

const Input = ({value, onChange, placeholder, icon, label, containerStyle}) => {
  const inputRef = useRef();

  return (
    <Box style={containerStyle}>
      {label && (
        <p className={styles.label}>{label}</p>
      )}

      <div
        className={styles.inputContainer}
        onClick={() => inputRef.current?.focus()}
      >
        {icon && <img
          src={icon}
          alt={'input icon'}
        />}

        <input
          ref={inputRef}
          className={styles.input}
          value={value}
          onChange={event => onChange(event.target.value)}
          type="text"
          placeholder={placeholder}
        />
      </div>
    </Box>
  );
};

export default Input;