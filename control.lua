 -- This file stands for console logging for FactoCord-3.0 integration
 -- Please configure as needed, any discord message will be sent in
 --  raw format if it starts with `0000-00-00 00:00:00 [DISCORD] `
 -- For more information visit https://github.com/maxsupermanhd/FactoCord-3.0
 -- If you have any question or comments join our Discord https://discord.gg/uNhtRH8

local FactoCordIntegration = {}

function FactoCordIntegration.PrintToDiscord(msg)
	localised_print({"", "0000-00-00 00:00:00 [DISCORD] ", msg})
end

script.on_event(defines.events.on_player_joined_game, function(event)
	local p = game.players[event.player_index];
	FactoCordIntegration.PrintToDiscord("**" .. p.name .. "** joined.");
	if(p.admin == true) then
		p.print("Welcome admin " .. p.name .. " to server!");
	else
		p.print("Welcome " .. p.name .. " to server!");
	end
end)

script.on_event(defines.events.on_player_left_game, function(event)
	local p = game.players[event.player_index];
	FactoCordIntegration.PrintToDiscord("**" .. p.name .. "** left.");
end)

script.on_event({defines.events.on_console_chat},
	function (e)
		if not e.player_index then
			return
		end
		if game.players[e.player_index].tag == "" then
			if game.players[e.player_index].admin then
				FactoCordIntegration.PrintToDiscord('(Admin) <' .. game.players[e.player_index].name .. '> ' .. e.message)
			else
				FactoCordIntegration.PrintToDiscord('<' .. game.players[e.player_index].name .. '> ' .. e.message)
			end
		else
			if game.players[e.player_index].admin then
				FactoCordIntegration.PrintToDiscord('(Admin) <' .. game.players[e.player_index].name .. '> ' .. game.players[e.player_index].tag .. " " .. e.message)
			else
				FactoCordIntegration.PrintToDiscord('<' .. game.players[e.player_index].name .. '> ' .. game.players[e.player_index].tag .. " " .. e.message)
			end
		end
	end
)


script.on_event(defines.events.on_player_died, function(event)
	local p = game.players[event.player_index];
	local c = event.cause
	if not c then
		FactoCordIntegration.PrintToDiscord("**" .. p.name .. "** died.");
	else
		local name = "Unknown";
		if c.type == "character" then
			name = c.player.name;
		elseif c.type == "spider-vehicle" then
			if c.entity_label then
				name = {"", c.localised_name, " " , c.entity_label};
			else
				name = {"", "a ", c.localised_name};
			end
		elseif c.type == "locomotive" then
			name = {"", c.localised_name, " " , c.backer_name};
		else
			name = {"", "a ", c.localised_name};
		end
		FactoCordIntegration.PrintToDiscord({"", "**", p.name, "** was killed by ", name, "."});
	end
end)
script.on_event(defines.events.on_player_kicked, function(event)
	local p = game.players[event.player_index];
	FactoCordIntegration.PrintToDiscord("**" .. p.name .. "** kicked.");
end)
script.on_event(defines.events.on_player_unbanned, function(event)
	FactoCordIntegration.PrintToDiscord("**" .. event.player_name .. "** unbanned.");
end)
script.on_event(defines.events.on_player_unmuted, function(event)
	local p = game.players[event.player_index];
	FactoCordIntegration.PrintToDiscord("**" .. p.name .. "** unmuted.");
end)
script.on_event(defines.events.on_player_banned, function(event)
	FactoCordIntegration.PrintToDiscord("**" .. event.player_name .. "** banned.");
end)
script.on_event(defines.events.on_player_muted, function(event)
	local p = game.players[event.player_index];
	FactoCordIntegration.PrintToDiscord("**" .. p.name .. "** muted.");
end)


return FactoCordIntegration;
