import { Component, createSignal, createResource, For } from "solid-js";
import TimePicker from "@rnwonder/solid-date-picker/timePicker";
import DatePicker, {
  PickerValue,
  TimeValue,
} from '@rnwonder/solid-date-picker';

import { CalendarIcon } from '../../Icons/Calendar';
import { NextIcon } from '../../Icons/Next';
import { PrevIcon } from '../../Icons/Prev';

import AgGridSolid from 'ag-grid-solid';
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { GridOptions } from 'ag-grid-community';

import { IDevice } from '../../types'
import { apiURL, apiGetDevices, apiGetDevice, apiGetDeviceEvents } from '../../constants'

interface DevicesProps {
    fetchAPI: (apiGetDevices: string) => any
    loadDevice: (device: IDevice) => void
    deviceDetail: () => any
    loadDeviceEvents: () => void
    toggleDeviceEvents: () => void
    area: () => string
    //deviceEventsToggle: () => any
    toggleReadDeviceEvents: () => number
    readDeviceDetail: () => IDevice
}

const Devices: Component<DevicesProps> = (props) => {


    function buildDateTime(fromOrtoOrCurrent: string): string {

        function fixZeroDateTime(dateTime: string): string {
          let receivedDateTime = new Date(dateTime)
          let hours = receivedDateTime.getHours()
          let minutes = receivedDateTime.getMinutes()
          let seconds = receivedDateTime.getSeconds()

          let currenthours = ("0" + hours).slice(-2);
          let currentminutes = ("0" + minutes).slice(-2);
          let currentseconds = ("0" + seconds).slice(-2);

          let returnFixedZeroDateTime = receivedDateTime.toLocaleDateString('sv-SE') + "T" + currenthours + ":" + currentminutes + ":" + currentseconds + ".000"
          //let returnFixedZeroDateTime = receivedDateTime.toLocaleDateString('sv-SE')+"T"+currenthours+":"+currentminutes+":"+currentseconds+".000Z"

          return returnFixedZeroDateTime
        }

        if (fromOrtoOrCurrent == 'from') {
          let fromDateTime = singleDateFrom().value?.selected?.split('T')[0] + 'T' + timeValueFrom().value.hour + ':' + timeValueFrom().value.minute + ':' + timeValueFrom().value.second
          console.log('fromDateTime: ', fromDateTime)
          //let fromDateTime = singleDateFrom().value.selected.split('T')[0] + 'T' + timeValueFrom().value.hour + ':' + timeValueFrom().value.minute + ':' + timeValueFrom().value.second
          return fromDateTime
        } else if (fromOrtoOrCurrent == 'to') {
          let toDateTime = singleDateTo().value?.selected?.split('T')[0] + 'T' + timeValueTo().value.hour + ':' + timeValueTo().value.minute + ':' + timeValueTo().value.second
          console.log('toDateTime: ', toDateTime)
          return toDateTime
        } else if (fromOrtoOrCurrent == 'current') {
          let currentDateTime = new Date()
          let returnCurrentDateTime = fixZeroDateTime(currentDateTime.toLocaleString())
          console.log('returnCurrentDateTime: ', returnCurrentDateTime)
          return returnCurrentDateTime
        } else if (fromOrtoOrCurrent == 'initFromDate') {
          let initFromDateTime = new Date()
          initFromDateTime.setHours(initFromDateTime.getHours() - 8)

          let initFromDateTimeStr = fixZeroDateTime(initFromDateTime.toLocaleString())
          console.log('initFromDateTime: ', initFromDateTimeStr)
          return initFromDateTimeStr
        } else {
          let initFromHour = new Date()
          initFromHour.setHours(initFromHour.getHours() - 8)
          let returnInitFromHour = initFromHour.getHours()
          console.log('initFromHour: ', returnInitFromHour)
          return returnInitFromHour.toString()
        }
    }

    const [deviceEventsToggle] = createResource(props.toggleReadDeviceEvents, async () => await props.fetchAPI(apiGetDeviceEvents.replace('{id}', props.readDeviceDetail().Id).replace('{start}', buildDateTime('from')).replace('{end}', buildDateTime('to'))));

    const [singleDateFrom, setSingleDateFrom] = createSignal<PickerValue>({
        label: (new Date(buildDateTime("initFromDate")).toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
        value: {
          selected: buildDateTime("initFromDate")
        },
      });
    const [ldeDF] = createResource(singleDateFrom, props.loadDeviceEvents)

    const [timeValueFrom, setTimeValueFrom] = createSignal<TimeValue>({
    value: {
        hour: parseInt(buildDateTime("initFromHour")),
        minute: new Date().getMinutes(),
        second: new Date().getSeconds()
    },
    label: (new Date(buildDateTime("initFromDate")).toLocaleTimeString()),
    });
    const [ldeTF] = createResource(timeValueFrom, props.loadDeviceEvents)

    const [singleDateTo, setSingleDateTo] = createSignal<PickerValue>({
    label: (new Date().toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
    value: {
        selected: buildDateTime("current")
    },
    });
    const [ldeDT] = createResource(singleDateTo, props.loadDeviceEvents)

    const [timeValueTo, setTimeValueTo] = createSignal<TimeValue>({
    value: {
        hour: new Date().getHours(),
        minute: new Date().getMinutes(),
        second: new Date().getSeconds()
    },
    label: (new Date().toLocaleTimeString()),
    });
    const [ldeTT] = createResource(timeValueTo, props.loadDeviceEvents)

    const columnDefsDevice = [
        { field: 'Id' },
        { field: 'name' },
        { field: 'label' },
        { field: 'type' },
        //{ field: 'room' },
      ];

      const columnDefsEvent = [
        { field: 'TimeStamp', width: 600 },
        { field: 'name' },
        { field: 'value' },
        { field: 'displayName' },
        { field: 'deviceId' },
        { field: 'descriptionText' },
        { field: 'unit' },
        { field: 'type' },
        { field: 'data' },
      ];

      const defaultColDef = {
        flex: 1,
        //minWidth: 100,
        filter: true,
        sortable: true,
        resizable: true,
      };

      const gridOptions: GridOptions = {
        autoSizeStrategy: {
          type: 'fitCellContents',
          //defaultMinWidth: 100,

          // columnLimits: [
          //   {
          //     colId: 'descriptionText',
          //     minWidth: 900
          //   }
          // ]
        },

        // other grid options ...
      }

    return (
        <>
            <div id='device_details' class={'w-full md:w-full bg-blue-200 h-[calc(100vh-160px)] text-black ' + (props.area() == 'device_details' ? '' : 'hidden')}>
                <div id="myGridDevice" class="content-center ag-theme-alpine " style={{ height: '100.0%' }}>
                {/* DeviceSelected: {deviceSelected()}
                Area: {area()} */}
                <AgGridSolid
                    rowData={props.deviceDetail()}
                    columnDefs={columnDefsDevice}
                    defaultColDef={defaultColDef}
                    animateRows={true}
                    autoSizeStrategy={gridOptions.autoSizeStrategy}
                    rowSelection='single'
                />
                </div>
            </div>

            <div id='device_events' class={'w-full md:w-full  bg-blue-200 h-[calc(100vh-160px)] text-black ' + (props.area() == 'device_events' ? '' : 'hidden')}>

                {/* ######################################################################################## */}
                {/* ############################### Calendar Component ##################################### */}
                <div id="date-time-menu" class='bg-red-500 h-[260px] flex flex-row flex-wrap p-3 space-x-2'>
                <div id="date-time-from" class="text-white flex flex-row flex-wrap p-3 space-x-2">
                    <p class={'text-lg '}>From:</p>
                    singleDateFrom: {singleDateFrom().value?.selected}
                    <br></br>
                    deviceEventsToggle: {deviceEventsToggle()}
                    <div>
                    <DatePicker
                        removeNavButtons
                        weekDaysType="single"
                        shouldHighlightWeekends
                        hideOutSideDays
                        afterNextButtonAreaJSX={({ handleNextMonth, handlePrevMonth }) => (
                        <div class="">
                            <button onClick={handlePrevMonth}>
                            <PrevIcon />
                            </button>
                            <button onClick={handleNextMonth}>
                            <NextIcon />
                            </button>
                        </div>
                        )}
                        renderInput={({ showDate, value }) => (
                        <div class="custom-input">
                            {/* <input value={value().label} readOnly /> */}
                            <input value={singleDateFrom().label} readOnly />
                            <button onClick={showDate}>
                            <CalendarIcon />
                            </button>
                        </div>
                        )}
                        onChange={() => { (props.toggleDeviceEvents()) }}
                        value={singleDateFrom}
                        setValue={setSingleDateFrom}
                        type={'single'}
                        monthSelectorFormat={'long'}
                    />
                    </div>
                    <div>
                    <TimePicker
                        value={timeValueFrom}
                        setValue={setTimeValueFrom}
                        shouldCloseOnSelect
                        allowedView={['hour', 'minute', 'second']}
                        renderInput={({ showTime, value }) => (
                        <div class="custom-input">
                            <input value={timeValueFrom().label} readOnly />
                            <button onClick={showTime}>
                            <CalendarIcon />
                            </button>
                        </div>
                        )}
                        timeAnalogClockCenterDotClass='text-red-500'
                        timeAnalogNumberClass='text-black'
                        timeAnalogClockHandClass='text-black'
                        timeAnalogWrapperClass='text-black'
                        timePickerBottomAreaClass='text-black'
                        timePickerMeridiemBtnClass='text-black'
                        timePickerTopAreaClass='text-black'
                        timePickerWrapperClass='text-black'
                    />
                    </div>
                </div>
                <div id="date-time-to" class="text-white flex flex-row flex-wrap p-3 space-x-2">
                    <p class={'text-lg '}>To:</p>
                    <div>
                    <DatePicker
                        removeNavButtons
                        weekDaysType="single"
                        shouldHighlightWeekends
                        hideOutSideDays
                        afterNextButtonAreaJSX={({ handleNextMonth, handlePrevMonth }) => (
                        <div class="">
                            <button onClick={handlePrevMonth}>
                            <PrevIcon />
                            </button>
                            <button onClick={handleNextMonth}>
                            <NextIcon />
                            </button>
                        </div>
                        )}
                        renderInput={({ showDate, value }) => (
                        <div class="custom-input">
                            <input value={singleDateTo().label} readOnly />
                            <button onClick={showDate}>
                            <CalendarIcon />
                            </button>
                        </div>
                        )}
                        onChange={() => { props.toggleDeviceEvents }}
                        value={singleDateTo}
                        setValue={setSingleDateTo}
                        type={'single'}
                        monthSelectorFormat={'long'}
                    />
                    </div>

                    <div>
                    <TimePicker
                        value={timeValueTo}
                        setValue={setTimeValueTo}
                        shouldCloseOnSelect
                        allowedView={['hour', 'minute', 'second']}
                        renderInput={({ showTime, value }) => (
                        <div class="custom-input">
                            <input value={timeValueTo().label} readOnly />
                            <button onClick={showTime}>
                            <CalendarIcon />
                            </button>
                        </div>
                        )}
                        onChange={() => { props.toggleDeviceEvents }}
                        timeAnalogClockCenterDotClass='text-red-500'
                        timeAnalogNumberClass='text-black'
                        timeAnalogClockHandClass='text-black'
                        timeAnalogWrapperClass='text-black'
                        timePickerBottomAreaClass='text-black'
                        timePickerMeridiemBtnClass='text-black'
                        timePickerTopAreaClass='text-black'
                        timePickerWrapperClass='text-black'
                    />
                    </div>
                </div>
                </div>

                <div id="myGridEvent" class="content-center ag-theme-alpine " style={{ height: '95.0%' }}>
                {(props.area() == 'device_events') && (
                    <AgGridSolid
                    rowData={deviceEventsToggle()}
                    columnDefs={columnDefsEvent}
                    defaultColDef={defaultColDef}
                    animateRows={true}
                    autoSizeStrategy={gridOptions.autoSizeStrategy}
                    />
                )}
                </div>
            </div>
        </>
    )
}

export default Devices;