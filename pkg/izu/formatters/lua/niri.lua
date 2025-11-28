local formatter = {}
local izu = izu

local capitalizations = {
	["super"] = "Super",
	["shift"] = "Shift",
	["ctrl"] = "Ctrl",
}

local function replace_capitalizations(keys)
	local output = {}
	for _, key in ipairs(keys) do
		local replacement = capitalizations[key]
		if replacement ~= nil then
			table.insert(output, replacement)
		else
			-- if length == 1, uppercase it
			if #key == 1 then
				table.insert(output, izu.uppercase(key))
			else
				table.insert(output, key)
			end
		end
	end
	return output
end

-- modifier order for `bind = Super+Shift, exec, echo hellow world
local modifiers = {
	"Super",
	"Shift",
	"Alt",
	"Ctrl",
	"Mod",
}

-- mousekeys for binds such as mouse:273, mouse:274, etc.
local mouse_keys = {
	["mouse_lmb"] = "MouseLeft",
	["mouse_rmb"] = "MouseRight",
	["mouse_mmb"] = "MouseMiddle",
}

local function replace_mousekey(key)
	if mouse_keys[key] then
		return mouse_keys[key]
	end

	if key:find("mouse_x") then
		local x = key:match("mouse_x(%d)")
		return "mouse:" .. (274 + tonumber(x))
	end

	return key
end

-- Formatter functions

function formatter.hotkey(args)
	return args.value[1] .. " { " .. args.value[2] .. " }"
end

function formatter.binding(args)
	if args.state == 1 then
		return table.concat(replace_capitalizations(args.value), "+")
	end
	return table.concat(args.value, "")
end

function formatter.multiple(args)
	return args.value
end

function formatter.single(args)
	local value = table.concat(args.value, "")
	if value == "_" then
		return
	end
	return { replace_mousekey(value) }
end

function formatter.string(args)
	return args.value
end

return formatter
