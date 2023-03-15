import React from 'react'
import {CheckboxOffIcon, CheckboxOnIcon} from "../../assets"
import {Box, Checkbox} from "@chakra-ui/react"
import styles from './CustomCheckbox.module.scss'
import {CheckboxProps} from "@chakra-ui/checkbox";

interface IProps extends CheckboxProps {
  label: string
  setChecked: (checked: boolean) => void
}

export const CustomCheckbox = ({label, checked, setChecked, ...checkboxProps}: IProps) => {
  const handleCheckboxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setChecked(event.target.checked)
  }

  const CustomIcon = () => {
    let src = checked ? CheckboxOnIcon : CheckboxOffIcon
    let alt = checked ? 'Checkbox On' : 'Checkbox Off'

    return <Box
      w={24}
      h={24}
      ml={12}
    >
      <img
        src={src}
        alt={alt}
      />
    </Box>
  }

  return (
    <Checkbox
      defaultChecked={checked}
      onChange={handleCheckboxChange}
      icon={<CustomIcon/>}
      spacing={0}
      checked={checked}
      className={styles.checkbox}
      {...checkboxProps}
    >
      {label}
    </Checkbox>
  )
}
