import type { Component } from 'solid-js';
import { onMount } from 'solid-js';
import { createSignal, createResource, For } from 'solid-js';
import TimePicker from "@rnwonder/solid-date-picker/timePicker";
import DatePicker, {
  PickerValue,
  TimeValue,
} from '@rnwonder/solid-date-picker';

import { CalendarIcon } from './Icons/Calendar';
import { NextIcon } from './Icons/Next';
import { PrevIcon } from './Icons/Prev';

import AgGridSolid from 'ag-grid-solid';
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { GridOptions } from 'ag-grid-community';

import { Sidebar, UpperMenu, Devices } from './components'
import { IDevice } from './types'

import { apiURL, apiGetDevices, apiGetDevice, apiGetDeviceEvents } from './constants'

const Device = {
  Id: '',
  name: '',
  label: '',
  type: '',
  room: ''
}

async function fetchAPI(apiUri: string) {
  console.log('fetchAPI: ', apiUri)
  try {
    const response = (await fetch(apiUri)).json();
    return response
  } catch (error) {
    console.error('Error fetching data:', error);
    throw error
  }
}

const App: Component = () => {
  const [area, setArea] = createSignal('device_details')

  const [readDeviceDetail, setReadDeviceDetail] = createSignal<IDevice>(Device);
  const [deviceDetail] = createResource(readDeviceDetail, async () => await fetchAPI(apiGetDevice.replace('{id}', readDeviceDetail().Id)));

  const [toggleReadDeviceEvents, settoggleReadDeviceEvents] = createSignal<number>(1);
  const toggleDeviceEvents = () => toggleReadDeviceEvents() == 1 ? settoggleReadDeviceEvents(2) : settoggleReadDeviceEvents(1)
  const [deviceEventsToggle] = createResource(toggleReadDeviceEvents, async () => await fetchAPI(apiGetDeviceEvents.replace('{id}', readDeviceDetail().Id).replace('{start}', buildDateTime('from')).replace('{end}', buildDateTime('to'))));

  const [singleDateFrom, setSingleDateFrom] = createSignal<PickerValue>({
    label: (new Date(buildDateTime("initFromDate")).toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
    value: {
      selected: buildDateTime("initFromDate")
    },
  });
  const [ldeDF] = createResource(singleDateFrom, loadDeviceEvents)

  const [timeValueFrom, setTimeValueFrom] = createSignal<TimeValue>({
    value: {
      hour: parseInt(buildDateTime("initFromHour")),
      minute: new Date().getMinutes(),
      second: new Date().getSeconds()
    },
    label: (new Date(buildDateTime("initFromDate")).toLocaleTimeString()),
  });
  const [ldeTF] = createResource(timeValueFrom, loadDeviceEvents)

  const [singleDateTo, setSingleDateTo] = createSignal<PickerValue>({
    label: (new Date().toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
    value: {
      selected: buildDateTime("current")
    },
  });
  const [ldeDT] = createResource(singleDateTo, loadDeviceEvents)

  const [timeValueTo, setTimeValueTo] = createSignal<TimeValue>({
    value: {
      hour: new Date().getHours(),
      minute: new Date().getMinutes(),
      second: new Date().getSeconds()
    },
    label: (new Date().toLocaleTimeString()),
  });
  const [ldeTT] = createResource(timeValueTo, loadDeviceEvents)
  const [deviceSelected, setDeviceSelected] = createSignal('false')

  function loadDevice(device: IDevice) {
    console.log('loadDevice: ', device)
    setDeviceSelected('true')
    console.log('deviceSelected: ', deviceSelected())
    setReadDeviceDetail(device)
    loadDeviceEvents()
  }

  function loadDeviceEvents() {
    console.log('loadDeviceEvents: ', readDeviceDetail)
    toggleDeviceEvents()
  }

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

  return (
    <>
      <header class="bg-slate-950 text-white p-3 text-2xl font-bold h-[60px]">HubitatDeviceEvents</header>
      <div class="flex flex-row h-[calc(100vh-100px)] bg-white">

        <Sidebar
          fetchAPI={fetchAPI}
          loadDevice={loadDevice}
        />

        <main id="main-area" class="h-full w-full overflow-auto flex flex-row flex-wrap p-2">

          <UpperMenu
            deviceSelected={deviceSelected}
            area={area}
            setArea={setArea}
          />

          <Devices
            fetchAPI={fetchAPI}
            loadDevice={loadDevice}
            area={area}
            deviceDetail={deviceDetail}
            //deviceEventsToggle={deviceEventsToggle}
            loadDeviceEvents={loadDeviceEvents}
            toggleDeviceEvents={toggleDeviceEvents}
            toggleReadDeviceEvents={toggleReadDeviceEvents}
            readDeviceDetail={readDeviceDetail}
          />

          <div id='device_details' class={'w-full md:w-full bg-blue-200 h-[calc(100vh-160px)] text-black ' + (area() == 'device_details' ? '' : 'hidden')}>
            {/* <p class='text-lg'>Device Details</p>
            <p class='text-lg'>Test readDeviceDetail() trigger: {readDeviceDetail().Id}</p>
            <p class='text-lg'>{readDeviceDetail() == Device ? '' : JSON.stringify(deviceDetail(), null, 2)}</p> */}

            <div id="myGridDevice" class="content-center ag-theme-alpine " style={{ height: '100.0%' }}>
              {/* DeviceSelected: {deviceSelected()}
              Area: {area()} */}
              <AgGridSolid
                rowData={deviceDetail()}
                columnDefs={columnDefsDevice}
                defaultColDef={defaultColDef}
                animateRows={true}
                autoSizeStrategy={gridOptions.autoSizeStrategy}
                rowSelection='single'
              />
            </div>

          </div>
          <div id='device_events' class={'w-full md:w-full  bg-blue-200 h-[calc(100vh-160px)] text-black ' + (area() == 'device_events' ? '' : 'hidden')}>

            {/* ######################################################################################## */}
            {/* ############################### Calendar Component ##################################### */}
            <div id="date-time-menu" class='bg-red-500 h-[60px] flex flex-row flex-wrap p-3 space-x-2'>
              <div id="date-time-from" class="text-white flex flex-row flex-wrap p-3 space-x-2">
                <p class={'text-lg '}>From:</p>
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
                    onChange={() => { toggleDeviceEvents }}
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
                    onChange={() => { toggleDeviceEvents }}
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
                    onChange={() => { toggleDeviceEvents }}
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
            {/* ######################################################################################## */}
            {/* <p class={'text-lg '}>Device Events:</p> */}
            {/* <p class='text-lg'>{readDeviceDetail() == Device ? '' : JSON.stringify(deviceEvents(), null, 2)}</p> */}
            {/* <p class='text-lg'>JSON Events: {JSON.stringify(deviceEventsToggle(), null, 2)}</p> */}
            {/* <p class='text-lg'>Events: {JSON.stringify(deviceEventsToggle(), null, 2)}</p> */}
            {/* <p class='text-lg'>{readTest()}</p> */}
            {/* <p class='text-lg'>Value of toggleReadDeviceEvents: {toggleReadDeviceEvents()}</p>
            <p class='text-lg'>Value of string toggleReadDeviceEvents: {toggleReadDeviceEvents().toString()}</p> */}

            <div id="myGridEvent" class="content-center ag-theme-alpine " style={{ height: '95.0%' }}>
            {/* DeviceSelected: {deviceSelected()}
            Area: {area()} */}
              {(area() == 'device_events') && (
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
        </main>
      </div >
      <footer class="h-[40px] text-white bg-black text-2xl px-2">footer</footer>
    </>
  );
};

export default App;
