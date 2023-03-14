import React, {useRef} from 'react';
import styles from './NumberInput.module.scss';
import {Box} from '@chakra-ui/react';
import CurrencyInput from 'react-currency-input-field';

const NumberInput = ({value, onChange, placeholder, icon, label, suffix, disabled}) => {
  const inputRef = useRef();

  const inputOnChange = (value) => {
    onChange(value);
  };

  return (
    <Box>
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

        <CurrencyInput
          ref={inputRef}
          className={styles.input}
          value={value}
          onValueChange={inputOnChange}
          placeholder={placeholder}
          suffix={suffix}
          disableGroupSeparators
          disabled={disabled}
          disableAbbreviations
          decimalsLimit={25}
        />
      </div>
    </Box>
  );
};

export default NumberInput;