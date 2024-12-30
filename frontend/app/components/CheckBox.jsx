import CheckIcon from "./icons/CheckIcon";

export default function CheckBox({title, value, id, onChanged}) {
  
    const classColors = 'border-gray-200 hover:border-blue-300 focus-active:border-blue-400 checked:bg-blue-500';
    const afterContent = "fill-white";
  
    function changeHandler(e) {
      onChanged(e);
    }
  
    return (
      <label className="flex flex-col">
        <div className="flex items-center">
          <div className="relative mr-2">
            {value && (
              <CheckIcon
                className={`${afterContent} absolute top-[20%] left-[15%] w-3.5 h-3.5`}
              />
            )}
            <input className={`
            appearance-none ${afterContent}
            flex justify-center items-center
            h-5 w-5 leading-5 text-sm outline-none border transition-colors rounded-md ${classColors}`}
              checked={value}
              onChange={changeHandler}
              type="checkbox"
              id={id}
            />
          </div>
          <span className={`text-sm transition-colors`}
          >
            {title}
          </span>
        </div>
      </label>
    )
  }