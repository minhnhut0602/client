// @flow
import Render from './index.render'
import type {DumbComponentMap} from '../../../constants/types/more'

const baseMock = {
  mode: 'codePageModeScanCode',
  textCode: 'derive jump shop federal member moral clip prevent vast wear critic purple mistake',
  qrCode: 'derive jump shop federal member moral clip prevent vast wear critic purple mistake',
  myDeviceRole: 'codePageDeviceRoleNewPhone',
  cameraBrokenMode: false,
  setCodePageMode: () => {},
  qrScanned: data => console.log('QR Scanned:', data),
  setCameraBrokenMode: () => {},
  textEntered: () => console.log('textEntered'),
  onChangeText: () => console.log('onChangeText'),
  onBack: () => console.log('onBack'),
  enterText: 'Foo Enter Text',
}

const scanCodeDeviceMock = {
  ...baseMock,
  mode: 'codePageModeScanCode',
  otherDeviceRole: 'codePageDeviceRoleExistingPhone',
}

const showTextDeviceMock = {
  ...baseMock,
  mode: 'codePageModeShowText',
  otherDeviceRole: 'codePageDeviceRoleExistingPhone',
}

const enterTextDeviceMock = {
  ...baseMock,
  mode: 'codePageModeEnterText',
  otherDeviceRole: 'codePageDeviceRoleExistingPhone',
}

const showCodeDeviceMock = {
  ...baseMock,
  mode: 'codePageModeShowCode',
  otherDeviceRole: 'codePageDeviceRoleExistingPhone',
}

const showCodeMock = {
  ...baseMock,
  mode: 'codePageModeShowCode',
  otherDeviceRole: 'codePageDeviceRoleExistingComputer',
}

const scanCodeMock = {
  ...baseMock,
  mode: 'codePageModeScanCode',
  otherDeviceRole: 'codePageDeviceRoleExistingComputer',
}

const showTextMock = {
  ...baseMock,
  mode: 'codePageModeShowText',
  otherDeviceRole: 'codePageDeviceRoleExistingComputer',
}

const enterTextMock = {
  ...baseMock,
  mode: 'codePageModeEnterText',
  otherDeviceRole: 'codePageDeviceRoleExistingComputer',
}

const dumbComponentMap: DumbComponentMap<Render> = {
  component: Render,
  mocks: {
    'Scan Code (Mobile)': scanCodeDeviceMock,
    'Show Text (Mobile)': showTextDeviceMock,
    'Enter Text (Mobile)': enterTextDeviceMock,
    'Show Code (Mobile)': showCodeDeviceMock,
    'Show Text': showTextMock,
    'Enter Text': enterTextMock,
    'Show Code': showCodeMock,
    'Scan Code': scanCodeMock,
  },
}

export default dumbComponentMap
