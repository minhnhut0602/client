// @flow
import * as CommonConstants from '../constants/common'
import {updateDebugConfig, updateReloading} from '../constants/dev'

import type {State} from '../constants/reducer'
import type {DebugConfig, DevAction} from '../constants/dev'

type DevState = {
  debugConfig: DebugConfig,
  hmrReloading: boolean,
}

const initialState: DevState = {
  debugConfig: {
    dumbFilter: '',
    dumbFullscreen: false,
    dumbIndex: 0,
  },
  hmrReloading: false,
}

export default function (state: DevState = initialState, action: DevAction): State {
  if (action.type === CommonConstants.resetStore) {
    return {...initialState}
  }

  if (action.type === updateDebugConfig) {
    return {
      ...state,
      debugConfig: {...state.debugConfig, ...action.payload},
    }
  }

  if (action.type === updateReloading && !action.error) {
    return {
      ...state,
      reloading: action.payload.reloading,
    }
  }
  return state
}
